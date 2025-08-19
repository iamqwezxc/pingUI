// main.go
package main

import (
	_ "github.com/lib/pq"

	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	"github.com/iamqwezxc/pingUI/Backend/pkg/handlers"
	"github.com/iamqwezxc/pingUI/Backend/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func main() {
	database.DBConnect(model.ConnStrUsers)
	oauth.InitOAuthConfig()

	r := gin.Default()

	// OAuth routes
	r.GET("/auth/google", handlers.GoogleLoginHandler)
	r.GET("/auth/google/callback", handlers.GoogleCallbackHandler)
	r.GET("/auth/yandex", handlers.YandexLoginHandler)
	r.GET("/auth/yandex/callback", handlers.YandexCallbackHandler)

	// Регистрируем остальные маршруты из website package
	// Если WBStarsWebSite регистрирует маршруты, нужно передать r
	registerWebsiteRoutes(r)

	// Или если WBStarsWebSite возвращает маршруты, то:
	// websiteRoutes := website.WBStarsWebSite()
	// r.Use(websiteRoutes)

	r.Run(":8080")
}

// Функция для регистрации маршрутов website
func registerWebsiteRoutes(r *gin.Engine) {
	// Здесь зарегистрируйте все маршруты из вашего website package
	// Например:
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to PingUI"})
	})

	// Добавьте другие маршруты по необходимости
}
