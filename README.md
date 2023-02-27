# How to use this package

1. Install the package
	
	`go get github.com/Zbyteio/zbyte-sso-verify-go`

2. Use the new keycloak access token verification as follows
```go
srv.GET("/a", func(ctx *gin.Context) {
  access_token := ctx.GetHeader("Authorization")
	isValid, err := middleware.MiddlewareHandler.VerifyOffline(access_token, "https://dplatdev.zbyte.io/kc")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error" : err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, isValid)
})
```