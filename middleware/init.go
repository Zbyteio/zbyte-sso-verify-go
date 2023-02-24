package middleware

var MiddlewareHandler middlewareInterface = &middlewareStruct{}

type middlewareStruct struct{}

type middlewareInterface interface {
	verifyOffline(accessToken string, keycloakDomain string) (jwtResponse *VerifyJwtTokenResponseKeycloak, errorData error)
}
