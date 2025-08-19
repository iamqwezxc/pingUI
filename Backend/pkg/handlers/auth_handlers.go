// handlers/auth_handlers.go
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	JSONJWT "github.com/iamqwezxc/pingUI/Backend/pkg/json_jwt"
	"github.com/iamqwezxc/pingUI/Backend/pkg/oauth"
	"golang.org/x/oauth2"
)

func GoogleLoginHandler(c *gin.Context) {
	log.Println("Initiating Google OAuth flow")
	url := oauth.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Redirecting to: %s", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallbackHandler(c *gin.Context) {
	log.Println("=== GOOGLE CALLBACK STARTED ===")

	// Логируем все параметры запроса
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")
	errorDescription := c.Query("error_description")

	log.Printf("Query parameters received:")
	log.Printf("  code: %s", code)
	log.Printf("  state: %s", state)
	log.Printf("  error: %s", errorParam)
	log.Printf("  error_description: %s", errorDescription)
	log.Printf("  Full URL: %s", c.Request.URL.String())

	// Проверяем ошибки OAuth
	if errorParam != "" {
		log.Printf("❌ OAuth error received: %s - %s", errorParam, errorDescription)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "OAuth authentication failed",
			"error_code":        errorParam,
			"error_description": errorDescription,
		})
		return
	}

	// Проверяем наличие authorization code
	if code == "" {
		log.Println("❌ No authorization code parameter in callback")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No authorization code received",
			"details": "The Google OAuth callback did not include an authorization code",
		})
		return
	}

	log.Println("✅ Authorization code received successfully")

	// Обмен кода на access token
	log.Println("Exchanging authorization code for access token...")
	token, err := oauth.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("❌ Token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange authorization code for token",
			"details": err.Error(),
		})
		return
	}
	log.Printf("✅ Token received successfully:")
	log.Printf("  Access Token: %s", token.AccessToken)
	log.Printf("  Token Type: %s", token.TokenType)
	log.Printf("  Expiry: %v", token.Expiry)
	if token.RefreshToken != "" {
		log.Printf("  Refresh Token: %s", token.RefreshToken)
	}

	// Получение информации о пользователе
	log.Println("Fetching user info from Google API...")
	userInfo, err := oauth.GetGoogleUserInfo(token)
	if err != nil {
		log.Printf("❌ Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve user information from Google",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ User info retrieved successfully:")
	log.Printf("  User ID: %s", userInfo.ID)
	log.Printf("  Email: %s", userInfo.Email)
	log.Printf("  Name: %s", userInfo.Name)
	log.Printf("  Picture: %s", userInfo.Picture)

	// Поиск или создание пользователя в базе данных
	log.Println("Checking if user exists in database...")
	user, err := database.DBFindUserByGoogleID(userInfo.ID)
	if err != nil {
		log.Printf("User not found by Google ID, creating new user...")

		newUser := model.User{
			FullName: userInfo.Name,
			Email:    userInfo.Email,
			GoogleID: userInfo.ID,
			Avatar:   userInfo.Picture,
			Role:     "user",
			Username: userInfo.Email, // Используем email как username
		}

		if err := database.DBCreateUserFromOAuth(&newUser); err != nil {
			log.Printf("❌ Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create user account",
				"details": err.Error(),
			})
			return
		}
		user = &newUser
		log.Printf("✅ New user created with ID: %d", user.ID)
	} else {
		log.Printf("✅ User found in database: ID %d", user.ID)
	}

	// Генерация JWT токена
	log.Println("Generating JWT token...")
	jwtToken, err := JSONJWT.GenerateToken(user.ID)
	if err != nil {
		log.Printf("❌ Failed to generate JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate authentication token",
			"details": err.Error(),
		})
		return
	}
	log.Printf("✅ JWT token generated successfully")

	// Успешный ответ
	log.Println("=== GOOGLE AUTHENTICATION COMPLETED SUCCESSFULLY ===")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authentication successful",
		"token":   jwtToken,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.FullName,
			"avatar":    user.Avatar,
			"role":      user.Role,
			"google_id": user.GoogleID,
		},
	})
}

func DebugOAuthConfig(c *gin.Context) {
	config := oauth.GoogleOAuthConfig

	debugInfo := gin.H{
		"client_id":    config.ClientID,
		"redirect_uri": config.RedirectURL,
		"scopes":       config.Scopes,
		"auth_url":     config.AuthCodeURL("debug-state"),
		"endpoint": gin.H{
			"auth_url":  config.Endpoint.AuthURL,
			"token_url": config.Endpoint.TokenURL,
		},
	}

	log.Printf("Debug OAuth Config: %+v", debugInfo)

	c.JSON(200, debugInfo)
}

// handlers/auth_handlers.go
func DebugAuthConfig(c *gin.Context) {
	config := oauth.GoogleOAuthConfig

	// Проверяем базовые настройки
	configCheck := gin.H{
		"client_id_exists":     config.ClientID != "",
		"client_secret_exists": config.ClientSecret != "",
		"redirect_uri":         config.RedirectURL,
		"scopes":               config.Scopes,
		"expected_redirect":    "http://localhost:3000/auth/google/callback",
		"matches_expected":     config.RedirectURL == "http://localhost:3000/auth/google/callback",
	}

	c.JSON(200, configCheck)
}
func TestSimpleAuth(c *gin.Context) {
	clientID := "727700104291-fsvh2ld6qss6st90rln3pagj53s7k0ud.apps.googleusercontent.com"
	redirectURI := "http://localhost:8080/auth/google/callback"

	url := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&access_type=offline&prompt=consent",
		clientID,
		redirectURI,
	)

	log.Printf("Simple auth URL: %s", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func DebugCallback(c *gin.Context) {
	log.Println("=== DEBUG CALLBACK ===")
	log.Printf("Full URL: %s", c.Request.URL.String())
	log.Printf("Query parameters: %+v", c.Request.URL.Query())

	// Логируем все query параметры
	for key, values := range c.Request.URL.Query() {
		log.Printf("  %s: %v", key, values)
	}

	c.JSON(200, gin.H{
		"url":          c.Request.URL.String(),
		"query_params": c.Request.URL.Query(),
	})
}
func YandexLoginHandler(c *gin.Context) {
	url := oauth.YandexOAuthConfig.AuthCodeURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func YandexCallbackHandler(c *gin.Context) {
	code := c.Query("code")

	token, err := oauth.YandexOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	userInfo, err := oauth.GetYandexUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Проверяем, есть ли пользователь в базе
	user, err := database.DBFindUserByYandexID(userInfo.ID)
	if err != nil {
		// Пользователь не найден, создаем нового
		avatarURL := fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/islands-200", userInfo.ProfileImage)

		newUser := model.User{
			FullName: userInfo.DisplayName,
			Email:    userInfo.Email,
			YandexID: userInfo.ID,
			Avatar:   avatarURL,
			Role:     "user",
		}

		if err := database.DBCreateUserFromOAuth(&newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		user = &newUser
	}

	// Генерируем JWT токен
	jwtToken, err := JSONJWT.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user":  user,
	})
}
