// pkg/wb_website/websites.go
package pkg

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"
	JSONJWT "github.com/iamqwezxc/pingUI/Backend/pkg/json_jwt"
)

func WBStarsWebSite(r *gin.Engine) {
	// ==================== USERS ====================
	// GET - получить всех пользователей
	r.GET("/users", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		users, err := database.GetSlice(db, "users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"users":   users,
			"count":   len(users),
		})
	})

	// GET - получить пользователя по ID
	r.GET("/users/:id", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		var user model.User
		err = db.QueryRow("SELECT * FROM users WHERE user_id = $1", userID).Scan(
			&user.ID, &user.FullName, &user.Username, &user.Email,
			&user.PasswordFirst, &user.PasswordSecond, &user.Role,
			&user.GoogleID, &user.YandexID, &user.Avatar,
		)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"user":    user,
		})
	})

	// POST - создать пользователя
	r.POST("/users", func(c *gin.Context) {
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
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "User created successfully",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		}
	})

	// PATCH - обновить пользователя
	r.PATCH("/users/:id", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var userUpdates model.User
		if err := c.ShouldBindJSON(&userUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		err = database.DBUpdateUserByID(userID, userUpdates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("User %d updated successfully", userID),
		})
	})

	// DELETE - удалить пользователя
	r.DELETE("/users/:id", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = database.DBDeleteUserByID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("User %d deleted successfully", userID),
		})
	})

	// ==================== COURSES ====================
	// GET - все курсы
	r.GET("/courses", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		rows, err := db.Query("SELECT * FROM courses")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var courses []model.Course
		for rows.Next() {
			var course model.Course
			err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.Thumbnail_url, &course.Instructor_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			courses = append(courses, course)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"courses": courses,
		})
	})

	// GET - курс по ID
	r.GET("/courses/:id", func(c *gin.Context) {
		courseID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		var course model.Course
		err = db.QueryRow("SELECT * FROM courses WHERE id = $1", courseID).Scan(
			&course.ID, &course.Title, &course.Description, &course.Thumbnail_url, &course.Instructor_id,
		)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"course":  course,
		})
	})

	// POST - создать курс
	r.POST("/courses", func(c *gin.Context) {
		course, err := JSONJWT.JSONtoStruct[model.Course](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataCourse(course)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Course created successfully",
		})
	})

	// PATCH - обновить курс
	r.PATCH("/courses/:id", func(c *gin.Context) {
		courseID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		var courseUpdates model.Course
		if err := c.ShouldBindJSON(&courseUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Здесь должна быть функция обновления курса
		// database.DBUpdateCourseByID(courseID, courseUpdates)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Course %d update endpoint - implement DBUpdateCourseByID", courseID),
		})
	})

	// DELETE - удалить курс
	r.DELETE("/courses/:id", func(c *gin.Context) {
		courseID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		// Здесь должна быть функция удаления курса
		// database.DBDeleteCourseByID(courseID)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Course %d delete endpoint - implement DBDeleteCourseByID", courseID),
		})
	})

	// ==================== LESSONS ====================
	// GET - все уроки
	r.GET("/lessons", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		rows, err := db.Query("SELECT * FROM lessons")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var lessons []model.Lesson
		for rows.Next() {
			var lesson model.Lesson
			err := rows.Scan(&lesson.ID, &lesson.Course_id, &lesson.Title, &lesson.Content, &lesson.Video_url, &lesson.Lesson_order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			lessons = append(lessons, lesson)
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"lessons": lessons,
		})
	})

	// GET - урок по ID
	r.GET("/lessons/:id", func(c *gin.Context) {
		lessonID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
			return
		}

		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		var lesson model.Lesson
		err = db.QueryRow("SELECT * FROM lessons WHERE lesson_id = $1", lessonID).Scan(
			&lesson.ID, &lesson.Course_id, &lesson.Title, &lesson.Content, &lesson.Video_url, &lesson.Lesson_order,
		)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"lesson":  lesson,
		})
	})

	// POST - создать урок
	r.POST("/lessons", func(c *gin.Context) {
		lesson, err := JSONJWT.JSONtoStruct[model.Lesson](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataLesson(lesson)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Lesson created successfully",
		})
	})

	// ==================== MATERIALS ====================
	// GET - все материалы
	r.GET("/materials", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		rows, err := db.Query("SELECT * FROM materials")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var materials []model.Material
		for rows.Next() {
			var material model.Material
			err := rows.Scan(&material.ID, &material.Lesson_id, &material.Title, &material.File_url, &material.TypeOfMaterial)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			materials = append(materials, material)
		}

		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"materials": materials,
		})
	})

	// POST - создать материал
	r.POST("/materials", func(c *gin.Context) {
		material, err := JSONJWT.JSONtoStruct[model.Material](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataMaterial(material)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Material created successfully",
		})
	})

	// ==================== ENROLLMENTS ====================
	// GET - все записи на курсы
	r.GET("/enrollments", func(c *gin.Context) {
		db := database.DBConnect(model.ConnStrUsers)
		defer db.Close()

		rows, err := db.Query("SELECT * FROM enrollments")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var enrollments []model.Enrollment
		for rows.Next() {
			var enrollment model.Enrollment
			err := rows.Scan(&enrollment.ID, &enrollment.User_id, &enrollment.Course_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			enrollments = append(enrollments, enrollment)
		}

		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"enrollments": enrollments,
		})
	})

	// POST - записаться на курс
	r.POST("/enrollments", func(c *gin.Context) {
		enrollment, err := JSONJWT.JSONtoStruct[model.Enrollment](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		database.DBAddDataEnrollment(enrollment)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Enrollment created successfully",
		})
	})

	// ==================== ROOT ====================
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
			"routes": []string{
				"GET    /users",
				"GET    /users/:id",
				"POST   /users",
				"PATCH  /users/:id",
				"DELETE /users/:id",
				"GET    /courses",
				"GET    /courses/:id",
				"POST   /courses",
				"PATCH  /courses/:id",
				"DELETE /courses/:id",
				"GET    /lessons",
				"GET    /lessons/:id",
				"POST   /lessons",
				"GET    /materials",
				"POST   /materials",
				"GET    /enrollments",
				"POST   /enrollments",
				"POST   /api/bash/execute",
				"GET    /api/bash/health",
			},
		})
	})

	log.Println("All CRUD routes registered successfully")
}
