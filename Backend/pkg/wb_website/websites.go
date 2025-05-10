package websites

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	JSONJWT "github.com/iamqwezxc/pingUI/Backend/pkg/json_jwt"
)

func WBStarsWebSite() {
	r := gin.Default()

	r.GET("/regist", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)

		err := database.TakeTable(db, c, "users")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer db.Close()

	})

	r.POST("/regist", func(c *gin.Context) {
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
