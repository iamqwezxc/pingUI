// pkg/oauth/oauth.go
package oauth

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/yandex"
)

var (
	YandexOAuthConfig *oauth2.Config
)

func InitOAuthConfig() {
	clientID := os.Getenv("YANDEX_CLIENT_ID")
	clientSecret := os.Getenv("YANDEX_CLIENT_SECRET")

	log.Printf("Yandex OAuth Config - ClientID: %s, ClientSecret: %s", clientID, clientSecret)

	if clientID == "" || clientSecret == "" {
		log.Fatal("YANDEX_CLIENT_ID or YANDEX_CLIENT_SECRET is not set in environment variables")
	}

	YandexOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/yandex/callback",
		Scopes:       []string{"login:email", "login:info"},
		Endpoint:     yandex.Endpoint,
	}

	log.Printf("Yandex OAuth Config Loaded successfully")
	log.Printf("Redirect URL: %s", YandexOAuthConfig.RedirectURL)
	log.Printf("Auth URL: %s", YandexOAuthConfig.Endpoint.AuthURL)
	log.Printf("Token URL: %s", YandexOAuthConfig.Endpoint.TokenURL)
}

type YandexUserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Email        string `json:"default_email"`
	DisplayName  string `json:"display_name"`
	ProfileImage string `json:"default_avatar_id"`
}

func GetYandexUserInfo(token *oauth2.Token) (*YandexUserInfo, error) {
	client := YandexOAuthConfig.Client(nil, token)
	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Логируем статус ответа
	log.Printf("Yandex API response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Логируем сырой ответ для отладки
	log.Printf("Yandex API raw response: %s", string(body))

	var userInfo YandexUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
