package oauth

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthService struct {
	config *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewGoogleOAuthService(clientID, clientSecret, redirectURL string) *GoogleOAuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &GoogleOAuthService{config: config}
}

func (g *GoogleOAuthService) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *GoogleOAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

func (g *GoogleOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (g *GoogleOAuthService) VerifyState(state, expected string) bool {
	return state == expected
}
