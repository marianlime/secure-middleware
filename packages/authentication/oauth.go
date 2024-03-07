package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)


func AdaptableConfigOAuth() *http.Client {
	ctx := context.Background()
	

		clientID := os.Getenv("OAUTH_CLIENT_ID")
		clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
		tokenURL := os.Getenv("OAUTH_TOKEN_URL")
		scopes := parseScopes(os.Getenv("OAUTH_SCOPES"))
	

	oauthConfig := &clientcredentials.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		TokenURL: tokenURL,
		Scopes: scopes,
	}

	return oauthConfig.Client(ctx)
}

func OauthMiddleware(next http.Handler) http.Handler{
	client := AdaptableConfigOAuth()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := client.Transport.(*oauth2.Transport).Source.Token()
		if err != nil {
			log.Printf("Failed to get OAuth token : %v", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		reqWithAuth := r.Clone(context.Background())
		reqWithAuth.Header.Set("Authorization", "Bearer"+token.AccessToken)

		next.ServeHTTP(w, reqWithAuth)
	})

}

func parseScopes(scopes string) []string {
	if scopes == "" {
		return nil
	}
	return strings.Split(scopes, ",")
}
