package middleware

import "github.com/dgrijalva/jwt-go"

type VerifyJwtTokenResponseKeycloak struct {
	Status bool                   `json:"status"`
	Header map[string]interface{} `json:"tokenHeader"`
	Data   jwt.Claims             `json:"data"`
}
type IntrospectTokenRequest struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
}
type VerifyJwtOnlineResponseKeycloak struct {
	Exp               int      `json:"exp,omitempty"`
	Iat               string   `json:"iat,omitempty"`
	AuthTime          string   `json:"auth_time,omitempty"`
	Jti               string   `json:"jti,omitempty"`
	Iss               string   `json:"iss,omitempty"`
	Aud               string   `json:"aud,omitempty"`
	Sub               string   `json:"sub,omitempty"`
	Typ               string   `json:"typ,omitempty"`
	Azp               int      `json:"azp,omitempty"`
	SessionState      string   `json:"session_state,omitempty"`
	Name              string   `json:"name,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
	FamilyName        string   `json:"family_name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Email             string   `json:"email,omitempty"`
	EmailVerified     string   `json:"email_verified,omitempty"`
	Address           struct{} `json:"address,omitempty"`
	Acr               string   `json:"acr,omitempty"`
	AllowedOrigins    []string `json:"allowed-origins,omitempty"`
	RealmAccess       struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"arealm_access,omitempty"`
	ResourceAccess struct {
		Account struct {
			Roles []string `json:"roles,omitempty"`
		}
	} `json:"resource_access,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Sid      string `json:"sid,omitempty"`
	ClientId string `json:"client_id,omitempty"`
	Username string `json:"username,omitempty"`
	Active   string `json:"active,omitempty"`
}
