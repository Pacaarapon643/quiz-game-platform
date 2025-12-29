package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quiz-game-backend/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type OAuthProvider string

const (
	ProviderGoogle   OAuthProvider = "google"
	ProviderFacebook OAuthProvider = "facebook"
)

type OAuthService struct {
	facebookConfig *oauth2.Config
	googleConfig   *oauth2.Config
}

type FacebookUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewOAuthService(cfg *config.OAuthConfig) *OAuthService {
	return &OAuthService{
		facebookConfig: &oauth2.Config{
			ClientID:     cfg.FacebookAppID,
			ClientSecret: cfg.FacebookAppSecret,
			RedirectURL:  cfg.FacebookRedirectURL,
			Scopes:       []string{"email", "public_profile"}, // Facebook ใช้ public_profile
			Endpoint:     facebook.Endpoint,
		},
		googleConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			}, // Google ใช้ full URL
			Endpoint: google.Endpoint,
		},
	}
}

// GetAuthURL สร้าง URL สำหรับ OAuth
func (s *OAuthService) GetAuthURL(provider OAuthProvider, state string) string {
	switch provider {
	case ProviderFacebook:
		return s.facebookConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	case ProviderGoogle:
		return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	default:
		return ""
	}
}

// ExchangeCode แลก code เป็น access token
func (s *OAuthService) ExchangeCode(ctx context.Context, provider OAuthProvider, code string) (*oauth2.Token, error) {
	switch provider {
	case ProviderFacebook:
		return s.facebookConfig.Exchange(ctx, code)
	case ProviderGoogle:
		return s.googleConfig.Exchange(ctx, code)
	default:
		return nil, fmt.Errorf("unsupported provider")
	}
}

// GetFacebookUserInfo ดึงข้อมูลจาก Facebook
func (s *OAuthService) GetFacebookUserInfo(ctx context.Context, token *oauth2.Token) (*FacebookUser, error) {
	client := s.facebookConfig.Client(ctx, token)

	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("facebook API error: %s", body)
	}

	var fbUser FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&fbUser); err != nil {
		return nil, err
	}

	return &fbUser, nil
}

// GetGoogleUserInfo ดึงข้อมูลจาก Google
func (s *OAuthService) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := s.googleConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API error: %s", body)
	}
	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}
	return &googleUser, nil
}
