
package auth

import (
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/golang-jwt/jwt"
	"errors"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(role string) (string, error){
	secretKey := os.Getenv("secret_key")
	if secretKey == "" {
		return "", errors.New("secret key for JWT is not set")
	}
	claims := CustomClaims{
		Role: role, 
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
func RoleMiddleware(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer")
			secretKey := os.Getenv("secret_key")
			if secretKey == "" {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error){
				return []byte(secretKey), nil
			})
			if err != nil{
				http.Error(w, "Unautharized", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
				if claims.Role == requiredRole{
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}