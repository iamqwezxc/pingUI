package websites

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	JSONJWT "github.com/iamqwezxc/pingUI/Backend/pkg/json_jwt"
)

func WBStarsWebSite() {
	r := gin.Default()

	r.GET("/regist", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		rows, err := db.Query("SELECT user_id, full_name, Username, Email, Password_Hash, Role FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var user_id int
			var full_name string
			var Username string
			var Email string
			var Password_Hash string
			var Role string
			if err := rows.Scan(&user_id, &full_name, &Username, &Email, &Password_Hash, &Role); err != nil {
				log.Fatal(err)
			}
			c.String(http.StatusOK, strconv.Itoa(user_id)+" ")
			c.String(http.StatusOK, full_name+" ")
			c.String(http.StatusOK, Username+" ")
			c.String(http.StatusOK, Email+" ")
			c.String(http.StatusOK, Password_Hash+" ")
			c.String(http.StatusOK, Role+" ")
			c.String(http.StatusOK, "\n")

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
