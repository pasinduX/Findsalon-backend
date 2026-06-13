package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"findsalon-backend/integrations"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     integrations.GoogleClientId,
		ClientSecret: integrations.GoogleClientSecret,
		RedirectURL:  integrations.GoogleRedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

func GetGoogleAuthURL(state string) string {
	return getGoogleOAuthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func ExchangeGoogleCode(code string) (*oauth2.Token, error) {
	token, err := getGoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func GetGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := getGoogleOAuthConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var info GoogleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
