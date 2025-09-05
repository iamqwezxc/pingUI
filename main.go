// main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	"github.com/iamqwezxc/pingUI/Backend/pkg/handlers"
	"github.com/iamqwezxc/pingUI/Backend/pkg/oauth"
	wb "github.com/iamqwezxc/pingUI/Backend/pkg/wb_website"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load("file.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Проверка обязательных переменных
	if os.Getenv("YANDEX_CLIENT_ID") == "" {
		log.Fatal("YANDEX_CLIENT_ID is not set")
	}
	if os.Getenv("YANDEX_CLIENT_SECRET") == "" {
		log.Fatal("YANDEX_CLIENT_SECRET is not set")
	}

	database.DBConnect(model.ConnStrUsers)
	oauth.InitOAuthConfig()

	r := gin.Default()

	// Yandex OAuth routes
	r.GET("/auth/yandex", handlers.YandexLoginHandler)
	r.GET("/auth/yandex/callback", handlers.YandexCallbackHandler)

	// Debug endpoint
	r.GET("/debug/yandex", handlers.DebugYandexConfig)

	// Website routes
	wb.WBStarsWebSite(r)

	r.POST("/api/bash/execute", handlers.BashExecuteHandler)
	r.GET("/api/bash/health", handlers.BashHealthHandler)

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
