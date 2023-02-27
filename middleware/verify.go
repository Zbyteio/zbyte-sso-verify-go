package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/Zbyteio/zbyte-sso-verify-go/config"
	"github.com/golang-jwt/jwt/v4"
)

func (m *middlewareStruct) VerifyOffline(accessToken string, baseUrl string, realm string) (jwtResponse *VerifyJwtOfflineTokenResponse, errorData error) {
	var errorMsg string = ""
	if accessToken == "" {
		errorMsg = "cannot get a valid access token"
		return nil, errors.New(errorMsg)
	}

	if baseUrl == "" {
		errorMsg = "cannot get a valid base url"
		return nil, errors.New(errorMsg)
	}

	if strings.HasPrefix(strings.ToLower(accessToken)[:7], "bearer") {
		accessToken = accessToken[7:]
	}

	jwks_url := fmt.Sprintf("%s/%s/%s/%s", baseUrl, config.REALMS, realm, config.CERTS_URL)

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
	jwks, err := keyfunc.Get(jwks_url, options)
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
	data := &VerifyJwtOfflineTokenResponse{
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

func (m *middlewareStruct) VerifyOnline(accessToken string, baseUrl string, realm string, clientId string, clientSecret string) (jwtResponse *VerifyJwtOnlineResponseKeycloak, errorData error) {
	var errorMsg string = ""

	//check if accesstoken passed is not empty
	if accessToken == "" {
		errorMsg = "cannot get a valid access token"
		return nil, errors.New(errorMsg)
	}

	//check if clientId passed is not empty
	if clientId == "" {
		errorMsg = "cannot get a valid client ID"
		return nil, errors.New(errorMsg)
	}

	//check if cleintSecret passed is not empty
	if clientSecret == "" {
		errorMsg = "cannot get a valid client secret"
		return nil, errors.New(errorMsg)
	}

	//check if baseUrl passed is not empty
	if baseUrl == "" {
		errorMsg = "cannot get a valid base url"
		return nil, errors.New(errorMsg)
	}

	//Check if accessToken containe "Bearer" at start, if so remove it to verify
	if strings.HasPrefix(strings.ToLower(accessToken)[:7], "bearer") {
		accessToken = accessToken[7:]
	}

	//Create a http client to make POST request
	client := &http.Client{}

	//form introspect URL using baseUrl passed
	introspect_url := fmt.Sprintf("%s/%s/%s/%s", baseUrl, config.REALMS, realm, config.INTROSPECT_URL)

	//form request body using params passed
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("token", accessToken)

	r, err := http.NewRequest(http.MethodPost, introspect_url, strings.NewReader(data.Encode()))
	if err != nil {
		errorMsg = "unable to create token introspect request"
		return nil, errors.New(errorMsg)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("url %s\n", introspect_url)

	res, err := client.Do(r)
	fmt.Printf("err %s\n", err)
	if err != nil {
		errorMsg = "failed API call"
		return nil, errors.New(errorMsg)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errorMsg = "unable to read json response"
		return nil, errors.New(errorMsg)
	}

	//convert API response to expected struct variable
	respBody := VerifyJwtOnlineResponseKeycloak{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		errorMsg = "unable to destructure json response"
		return nil, errors.New(errorMsg)
	}
	return &respBody, nil
}
