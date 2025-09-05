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

func YandexLoginHandler(c *gin.Context) {
	url := oauth.YandexOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Yandex OAuth URL: %s", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func YandexCallbackHandler(c *gin.Context) {
	log.Println("=== YANDEX CALLBACK STARTED ===")

	// Логируем все параметры запроса
	for key, values := range c.Request.URL.Query() {
		log.Printf("Query param %s: %v", key, values)
	}

	code := c.Query("code")
	errorParam := c.Query("error")
	errorDescription := c.Query("error_description")

	if errorParam != "" {
		log.Printf("❌ OAuth error: %s - %s", errorParam, errorDescription)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":             "OAuth authentication failed",
			"error_code":        errorParam,
			"error_description": errorDescription,
		})
		return
	}

	if code == "" {
		log.Println("❌ No authorization code received")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code received"})
		return
	}

	log.Printf("✅ Authorization code received: %s", code)

	// Обмен кода на токен с дополнительной информацией об ошибке
	log.Println("Exchanging authorization code for access token...")
	token, err := oauth.YandexOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("❌ Token exchange failed: %v", err)

		// Попробуем получить больше информации об ошибке
		if oauthErr, ok := err.(*oauth2.RetrieveError); ok {
			log.Printf("OAuth error body: %s", string(oauthErr.Body))
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange token",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ Token received: AccessToken=%s, TokenType=%s", token.AccessToken, token.TokenType)

	log.Println("Fetching user info from Yandex API...")
	userInfo, err := oauth.GetYandexUserInfo(token)
	if err != nil {
		log.Printf("❌ Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user info",
			"details": err.Error(),
		})
		return
	}

	log.Printf("✅ User info retrieved: ID=%s, Email=%s, Name=%s", userInfo.ID, userInfo.Email, userInfo.DisplayName)

	// Проверяем, есть ли пользователь в базе
	user, err := database.DBFindUserByYandexID(userInfo.ID)
	if err != nil {
		log.Println("User not found, creating new user...")

		avatarURL := ""
		if userInfo.ProfileImage != "" {
			avatarURL = fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/islands-200", userInfo.ProfileImage)
		}

		// handlers/auth_handlers.go
		// В функции YandexCallbackHandler
		newUser := model.User{
			FullName: userInfo.DisplayName,
			Email:    userInfo.Email,
			YandexID: userInfo.ID,
			Avatar:   avatarURL,
			Role:     "student", // Указываем правильную роль
			Username: userInfo.Login,
		}

		if newUser.Username == "" {
			newUser.Username = userInfo.Email // Используем email если login пустой
		}

		if err := database.DBCreateUserFromOAuth(&newUser); err != nil {
			log.Printf("❌ Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create user",
				"details": err.Error(),
			})
			return
		}
		user = &newUser
		log.Printf("✅ New user created with ID: %d", user.ID)
	} else {
		log.Printf("✅ User found in database: ID %d", user.ID)
	}

	// Генерируем JWT токен
	jwtToken, err := JSONJWT.GenerateToken(user.ID)
	if err != nil {
		log.Printf("❌ Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"details": err.Error(),
		})
		return
	}

	log.Println("=== YANDEX AUTHENTICATION COMPLETED SUCCESSFULLY ===")
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
			"yandex_id": user.YandexID,
		},
	})
}

// Добавим debug endpoint для проверки конфигурации
func DebugYandexConfig(c *gin.Context) {
	config := oauth.YandexOAuthConfig

	debugInfo := gin.H{
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret != "", // только факт наличия, не значение
		"redirect_uri":  config.RedirectURL,
		"scopes":        config.Scopes,
		"auth_url":      config.Endpoint.AuthURL,
		"token_url":     config.Endpoint.TokenURL,
	}

	c.JSON(200, debugInfo)
}
