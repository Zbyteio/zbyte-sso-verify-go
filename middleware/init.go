package middleware

type middlewareStruct struct{}

var MiddlewareHandler middlewareInterface = &middlewareStruct{}

type middlewareInterface interface {
	VerifyOffline(accessToken string, baseUrl string, realm string) (jwtResponse *VerifyJwtOfflineTokenResponse, errorData error)
	VerifyOnline(accessToken string, baseUrl string, realm string, clientId string, clientSecret string) (jwtResponse *VerifyJwtOnlineResponseKeycloak, errorData error)
}
