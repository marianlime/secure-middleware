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
func RoleMiddleware(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
				tokenString := r.Header.Get("Authorization")
				tokenString = strings.TrimPrefix(tokenString, "Bearer")
			
				
				token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error){
					return []byte("Secret"), nil
				})


				if err != nil {
					http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
					return
				}

				if claims, ok := token.Claims.(*CustomClaims); ok && token.valid{
					if claims.Role == requiredRole {
						next.ServeHTTP(w, r)
						return
					}
				}

				http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}