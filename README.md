# How to use this package

1. Install the package
	
	`go get github.com/Zbyteio/zbyte-sso-verify-go`

2. Use the new access token verification as follows
```go
srv.GET("/exampleUsingVerifyOffline", func(ctx *gin.Context) {
  	access_token := ctx.GetHeader("Authorization")
  	baseUrl:= "https://dplatdev.zbyte.io/kc"
	realm:="community"

	isValid, err := middleware.MiddlewareHandler.VerifyOffline(access_token,baseUrl, realm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, isValid)
})


srv.GET("/exampleUsingVerifyOnline", func(ctx *gin.Context) {
  	access_token := ctx.GetHeader("Authorization")
  	baseUrl:="https://dplatdev.zbyte.io/kc"
  	client_id:= "lcnc-app"
  	client_secret:="8T4QEaE5aRIbn9bChDPyVIq4dnjJsxaW"
	realm:="community"

	isValid, err := middleware.MiddlewareHandler.VerifyOnline(access_token,baseUrl ,realm,clientId, client_secret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, isValid)
})
```