package auth

import (
	"net/http"
	"fmt"
)

type Authenticator interface {
	Authenticate(r* http.Request) (*AuthResult, error)
}

type AuthResult struct{
	User string
	Roles []string
	Authed bool
	Err error
}

type AuthManager struct {
	jwtAuth *JWTAuthenticator
	oauthAuth *OAuthAuthenticator
}

func newAuthManager() *AuthManager{
	return &AuthManager{
		jwtAuth: &JWTAuthentication{},
		oathAuth: &OAuthAutheticator{},
	}
}