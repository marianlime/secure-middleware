package auth

import(
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(role string) (string, error){
	claims := CustomClaims{
		Role: role, 
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("Secret"))
}