package websites

import (
	"net/http"
	database "studyproject/pkg/db_database"
	JSONJWT "studyproject/pkg/json_jwt"

	"github.com/gin-gonic/gin"
)

func WBStarsWebSite() {
	r := gin.Default()

	r.GET("/user", func(c *gin.Context) {
		c.String(http.StatusOK, "Главная страница")
	})

	r.POST("/user", func(c *gin.Context) {
		user, err := JSONJWT.JSONtoStruct(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.Password, err = JSONJWT.HashPassword(user.Password)
		database.DBAddDataUsers(user)
	})

	r.Run(":8080")

}
