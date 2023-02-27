package middleware

type middlewareStruct struct{}

var MiddlewareHandler middlewareInterface = &middlewareStruct{}

type middlewareInterface interface {
	VerifyOffline(accessToken string, baseUrl string) (jwtResponse *VerifyJwtOfflineTokenResponse, errorData error)
}
