package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

var MiddlewareHandler middlewareInterface = &middlewareStruct{}

type middlewareStruct struct{}

type middlewareInterface interface {
	verifyOffline(accessToken string, keycloakBaseUrl string) (jwtResponse *VerifyJwtTokenResponseKeycloak, errorData error)

	verifyOnline(accessToken string, keycloakBaseUrl string, clientId string, clientSecret string) (jwtResponse *VerifyJwtOnlineResponseKeycloak, errorData error)
}

func (m *middlewareStruct) verifyOffline(accessToken string, keycloakBaseUrl string) (jwtResponse *VerifyJwtTokenResponseKeycloak, errorData error) {
	var errorMsg string = ""
	if accessToken == "" {
		errorMsg = "cannot get a valid access token"
		return nil, errors.New(errorMsg)
	}

	if keycloakBaseUrl == "" {
		errorMsg = "cannot get a valid keycloak base url"
		return nil, errors.New(errorMsg)
	}

	if strings.HasPrefix(strings.ToLower(accessToken)[:7], "bearer") {
		accessToken = accessToken[7:]
	}

	keycloak_jwks_url := fmt.Sprintf("%s/realms/community/protocol/openid-connect/certs", keycloakBaseUrl)

	// Create a context that, when cancelled, ends the JWKS background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			errorMsg = fmt.Sprintf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
			log.Println(errorMsg)
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(keycloak_jwks_url, options)
	if err != nil {
		cancel()
		errorMsg = fmt.Sprintf("Failed to create JWKS from resource at the given URL.\nError: %s", err.Error())
		fmt.Println(errorMsg)
	}

	// Parse the JWT.
	token, err := jwt.Parse(accessToken, jwks.Keyfunc)
	if err != nil {
		cancel()
		errorMsg = fmt.Sprintf("Failed to parse the JWT.\nError: %s", err.Error())
		return nil, errors.New(errorMsg)
	}

	// check if any error produced
	if errorMsg != "" {
		cancel()
		return nil, errors.New(errorMsg)
	}

	// create response object to respond with
	data := &VerifyJwtTokenResponseKeycloak{
		Status: token.Valid,
		Header: token.Header,
		Data:   token.Claims,
	}

	// End the background refresh goroutine when it's no longer needed.
	cancel()

	// This will be ineffectual because the line above this canceled the parent context.Context.
	// This method call is idempotent similar to context.CancelFunc.
	jwks.EndBackground()

	// return data
	return data, nil
}

func (m *middlewareStruct) verifyOnline(accessToken string, keycloakBaseUrl string, clientId string, clientSecret string) (jwtResponse *VerifyJwtOnlineResponseKeycloak, errorData error) {
	var errorMsg string = ""
	if accessToken == "" {
		errorMsg = "cannot get a valid access token"
		return nil, errors.New(errorMsg)
	}

	if clientId == "" {
		errorMsg = "cannot get a valid client ID"
		return nil, errors.New(errorMsg)
	}

	if clientSecret == "" {
		errorMsg = "cannot get a valid client secret"
		return nil, errors.New(errorMsg)
	}

	if keycloakBaseUrl == "" {
		errorMsg = "cannot get a valid keycloak base url"
		return nil, errors.New(errorMsg)
	}

	if strings.HasPrefix(strings.ToLower(accessToken)[:7], "bearer") {
		accessToken = accessToken[7:]
	}

	// get user data from user-service
	// using the userId
	client := &http.Client{}

	keycloak_introspect_url := fmt.Sprintf("%s/realms/community/protocol/openid-connect/token/introspect", keycloakBaseUrl)

	postData := IntrospectTokenRequest{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	}
	json_data, _ := json.Marshal(postData)
	r, err := http.NewRequest(http.MethodPost, keycloak_introspect_url, bytes.NewBuffer(json_data))
	if err != nil {
		errorMsg = "unable to create token introspect request"
		return nil, errors.New(errorMsg)
	}
	r.Header.Add("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		errorMsg = "failed API call"
		return nil, errors.New(errorMsg)
	}
	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
	fmt.Println("response Headers:", res.Header)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errorMsg = "unable to read json response"
		return nil, errors.New(errorMsg)
	}
	var respBody VerifyJwtOnlineResponseKeycloak
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		errorMsg = "unable to destructure json response"
		return nil, errors.New(errorMsg)
	}
	fmt.Println("response Body:", string(body))
	return &respBody, nil
}
