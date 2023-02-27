package middleware

type middlewareStruct struct{}

var MiddlewareHandler middlewareInterface = &middlewareStruct{}

type middlewareInterface interface {
	VerifyOffline(accessToken string, keycloakBaseUrl string) (jwtResponse *VerifyJwtTokenResponseKeycloak, errorData error)
}
