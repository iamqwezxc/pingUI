// pkg/oauth/oauth.go
package oauth

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/yandex"
)

var (
	GoogleOAuthConfig *oauth2.Config
	YandexOAuthConfig *oauth2.Config
)

func InitOAuthConfig() {
	// Временно используем хардкод для точного совпадения
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     "727700104291-fsvh2ld6qss6st90rln3pagj53s7k0ud.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/google/callback", // порт 3000
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	log.Printf("OAuth Config Loaded: %s", GoogleOAuthConfig.RedirectURL)

	YandexOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("YANDEX_CLIENT_ID"),
		ClientSecret: os.Getenv("YANDEX_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/yandex/callback",
		Scopes:       []string{"login:email", "login:info"},
		Endpoint:     yandex.Endpoint,
	}
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type YandexUserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Email        string `json:"default_email"`
	DisplayName  string `json:"display_name"`
	ProfileImage string `json:"default_avatar_id"`
}

func GetGoogleUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := GoogleOAuthConfig.Client(nil, token)
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

func GetYandexUserInfo(token *oauth2.Token) (*YandexUserInfo, error) {
	client := YandexOAuthConfig.Client(nil, token)
	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo YandexUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
