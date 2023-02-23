package middleware

import "github.com/dgrijalva/jwt-go"

type VerifyJwtTokenResponseKeycloak struct {
	Status bool                   `json:"status"`
	Header map[string]interface{} `json:"tokenHeader"`
	Data   jwt.Claims             `json:"data"`
}
