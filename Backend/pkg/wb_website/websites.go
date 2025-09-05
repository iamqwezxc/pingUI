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

func WBStarsWebSite(r *gin.Engine) {
	// Маршруты для получения данных из таблиц
	r.GET("/regist", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		err := database.TakeTable(db, c, "users")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	r.GET("/courses", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		err := database.TakeTable(db, c, "courses")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	r.GET("/lessons", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		err := database.TakeTable(db, c, "lessons")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	r.GET("/materials", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		err := database.TakeTable(db, c, "materials")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	r.GET("/enrollments", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		err := database.TakeTable(db, c, "enrollments")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	})

	// Маршруты для добавления данных
	r.POST("/regist", func(c *gin.Context) {
		user, err := JSONJWT.JSONtoStruct[model.User](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.PasswordFirst == user.PasswordSecond {
			user.PasswordFirst, err = JSONJWT.HashPassword(user.PasswordFirst)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}

			database.DBAddDataUsers(user)
			c.JSON(http.StatusOK, gin.H{"success": true})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		}
	})

	r.POST("/courses", func(c *gin.Context) {
		course, err := JSONJWT.JSONtoStruct[model.Course](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataCourse(course)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.POST("/lessons", func(c *gin.Context) {
		lesson, err := JSONJWT.JSONtoStruct[model.Lesson](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataLesson(lesson)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.POST("/materials", func(c *gin.Context) {
		material, err := JSONJWT.JSONtoStruct[model.Material](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataMaterial(material)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.POST("/enrollments", func(c *gin.Context) {
		enrollment, err := JSONJWT.JSONtoStruct[model.Enrollment](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataEnrollment(enrollment)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Маршруты для обновления и удаления пользователей
	r.PUT("/users/edit/:id", func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid user ID format"})
			return
		}

		var userUpdates model.User
		if err := c.ShouldBindJSON(&userUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"Error":   "Invalid JSON payload: " + err.Error(),
			})
			return
		}

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

	// Маршрут для логина
	r.GET("/login", func(c *gin.Context) {
		c.String(http.StatusOK, "Логин")
	})

	r.POST("/login", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		c.String(http.StatusOK, fmt.Sprintf("%v", db))
	})

	// Корневой маршрут для проверки работы сервера
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	log.Println("Website routes registered successfully")
}
