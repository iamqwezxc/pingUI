package pkg

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	r.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "Логин")

	})

	r.POST("/login", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)

		c.String(http.StatusOK, fmt.Sprintf("%v", db))

		defer db.Close()

	})

	r.POST("/regist", func(c *gin.Context) {
		user, err := JSONJWT.JSONtoStruct[model.User](c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.PasswordFirst == user.PasswordSecond {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})

			user.PasswordFirst, err = JSONJWT.HashPassword(user.PasswordFirst)
			database.DBAddDataUsers(user)

		} else {
			c.JSON(http.StatusBadGateway, gin.H{
				"succes": false,
				"Error":  err.Error(),
			})

		}
	})
	r.PUT("/users/edit/:id", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var userUpdates model.User
		// Bind JSON without using JSONtoStruct to have more control over the response
		if err := c.ShouldBindJSON(&userUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"Error":   "Invalid JSON payload: " + err.Error(),
			})
			return
		}

		// The ID in model.User from JSON body is ignored; userID from path is authoritative.
		// userUpdates.ID is not used by DBUpdateUserByID for query condition.

		err = database.DBUpdateUserByID(userID, userUpdates)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found") {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "error": errMsg})
			} else if strings.Contains(errMsg, "no fields provided") {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": errMsg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update user: " + errMsg})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("User %d updated successfully", userID)})
	})

	// PUT request to delete a user by ID (using PUT method as requested)
	// Path: /users/delete/:id
	r.PUT("/users/delete/:id", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		err = database.DBDeleteUserByID(userID)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found") {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "error": errMsg})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete user: " + errMsg})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": fmt.Sprintf("User %d deleted successfully", userID)})
	})

	fmt.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
