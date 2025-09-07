package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
	JSONJWT "github.com/iamqwezxc/pingUI/Backend/pkg/json_jwt"
)

func DBConnect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	return db
}
func DBUpdateUserByID(userID int, userUpdates model.User) error {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var setClauses []string
	var args []interface{}
	argCounter := 1

	if userUpdates.FullName != "" {
		setClauses = append(setClauses, fmt.Sprintf("full_name = $%d", argCounter))
		args = append(args, userUpdates.FullName)
		argCounter++
	}
	if userUpdates.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("Username = $%d", argCounter))
		args = append(args, userUpdates.Username)
		argCounter++
	}
	if userUpdates.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("Email = $%d", argCounter))
		args = append(args, userUpdates.Email)
		argCounter++
	}
	if userUpdates.PasswordFirst != "" {
		hashedPassword, err := JSONJWT.HashPassword(userUpdates.PasswordFirst)
		if err != nil {
			log.Printf("Error hashing password for update: %v", err)
			return fmt.Errorf("failed to hash new password: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("Password_Hash = $%d", argCounter))
		args = append(args, hashedPassword)
		argCounter++
	}
	if userUpdates.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("Role = $%d", argCounter))
		args = append(args, userUpdates.Role)
		argCounter++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = $%d", strings.Join(setClauses, ", "), argCounter)
	args = append(args, userID)

	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating user %d: %v", userID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user %d update: %v", userID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found or no data changed", userID)
	}

	return nil
}

func DBDeleteUserByID(userID int) error {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	result, err := db.Exec("DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Error deleting user %d: %v", userID, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected for user %d deletion: %v", userID, err)

	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}
	return nil
}

func DBAddDataUsers(user model.User) {
	db := DBConnect(model.ConnStrUsers)

	// Устанавливаем роль по умолчанию если не указана
	role := user.Role
	if role == "" {
		role = "student"
	} else {
		// Проверяем валидность роли
		validRoles := map[string]bool{"student": true, "teacher": true, "admin": true}
		if !validRoles[role] {
			role = "student" // fallback to student if invalid
		}
	}

	passwordHash := user.PasswordFirst
	if passwordHash == "" {
		passwordHash = "oauth_user_no_password"
	}

	_, err := db.Exec(
		"INSERT INTO users (full_name, Username, Email, Password_Hash, Role) VALUES ($1, $2, $3, $4, $5)",
		user.FullName,
		user.Username,
		user.Email,
		passwordHash,
		role, // Используем проверенную роль
	)
	if err != nil {
		log.Printf("Error adding user: %v", err)
	}
	defer db.Close()
}

// Добавьте эти функции в pkg/db_database/databasepostgres.go

// DBUpdateCourseByID - обновление курса
func DBUpdateCourseByID(courseID int, courseUpdates model.Course) error {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	_, err := db.Exec(`
		UPDATE courses 
		SET title = $1, description = $2, thumbnail_url = $3, instructor_id = $4 
		WHERE id = $5`,
		courseUpdates.Title, courseUpdates.Description,
		courseUpdates.Thumbnail_url, courseUpdates.Instructor_id, courseID,
	)

	return err
}

// DBDeleteCourseByID - удаление курса
func DBDeleteCourseByID(courseID int) error {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	_, err := db.Exec("DELETE FROM courses WHERE id = $1", courseID)
	return err
}

// DBGetCourseByID - получение курса по ID
func DBGetCourseByID(courseID int) (*model.Course, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var course model.Course
	err := db.QueryRow("SELECT * FROM courses WHERE id = $1", courseID).Scan(
		&course.ID, &course.Title, &course.Description, &course.Thumbnail_url, &course.Instructor_id,
	)

	return &course, err
}
func DBAddDataCourse(course model.Course) {
	db := DBConnect(model.ConnStrUsers)
	_, err := db.Exec(
		"INSERT INTO courses (title, description, thumbnail_url, instructor_id) VALUES ($1, $2, $3, $4)",
		course.Title,
		course.Description,
		course.Thumbnail_url,
		course.Instructor_id,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func DBAddDataLesson(lesson model.Lesson) {
	db := DBConnect(model.ConnStrUsers)
	_, err := db.Exec(
		"INSERT INTO lessons (course_id, title, content, video_url, lesson_order) VALUES ($1, $2, $3, $4, $5)",
		lesson.Course_id,
		lesson.Title,
		lesson.Content,
		lesson.Video_url,
		lesson.Lesson_order,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func DBAddDataMaterial(material model.Material) {
	db := DBConnect(model.ConnStrUsers)
	_, err := db.Exec(
		"INSERT INTO materials (lesson_id, title, file_url, type) VALUES ($1, $2, $3, $4)",
		material.Lesson_id,
		material.Title,
		material.File_url,
		material.TypeOfMaterial,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func DBAddDataEnrollment(enrollment model.Enrollment) {
	db := DBConnect(model.ConnStrUsers)
	_, err := db.Exec(
		"INSERT INTO enrollments (user_id, course_id) VALUES ($1, $2)",
		enrollment.User_id,
		enrollment.Course_id,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func GetSlice(db *sql.DB, tableName string) ([]model.User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		// Сканируем ВСЕ 10 полей
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Username,
			&user.Email,
			&user.PasswordFirst,
			&user.PasswordSecond,
			&user.Role,
			&user.GoogleID,
			&user.YandexID,
			&user.Avatar,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// DBGetUserByID - получение пользователя по ID
func DBGetUserByID(userID int) (*model.User, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var user model.User
	err := db.QueryRow(`
        SELECT user_id, full_name, username, email, password_hash, role, google_id, yandex_id, avatar 
        FROM users WHERE user_id = $1`,
		userID,
	).Scan(
		&user.ID, &user.FullName, &user.Username, &user.Email,
		&user.PasswordFirst, &user.Role, &user.GoogleID, &user.YandexID, &user.Avatar,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DBGetAllCourses - получение всех курсов
func DBGetAllCourses() ([]model.Course, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.Thumbnail_url, &course.Instructor_id)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}

// DBGetAllLessons - получение всех уроков
func DBGetAllLessons() ([]model.Lesson, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM lessons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []model.Lesson
	for rows.Next() {
		var lesson model.Lesson
		err := rows.Scan(&lesson.ID, &lesson.Course_id, &lesson.Title, &lesson.Content, &lesson.Video_url, &lesson.Lesson_order)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// DBGetAllMaterials - получение всех материалов
func DBGetAllMaterials() ([]model.Material, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM materials")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []model.Material
	for rows.Next() {
		var material model.Material
		err := rows.Scan(&material.ID, &material.Lesson_id, &material.Title, &material.File_url, &material.TypeOfMaterial)
		if err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}

	return materials, nil
}

// DBGetAllEnrollments - получение всех записей
func DBGetAllEnrollments() ([]model.Enrollment, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM enrollments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []model.Enrollment
	for rows.Next() {
		var enrollment model.Enrollment
		err := rows.Scan(&enrollment.ID, &enrollment.User_id, &enrollment.Course_id)
		if err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

// DBGetLessonByID - получение урока по ID
func DBGetLessonByID(lessonID int) (*model.Lesson, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var lesson model.Lesson
	err := db.QueryRow("SELECT * FROM lessons WHERE lesson_id = $1", lessonID).Scan(
		&lesson.ID, &lesson.Course_id, &lesson.Title, &lesson.Content, &lesson.Video_url, &lesson.Lesson_order,
	)

	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// DBGetMaterialByID - получение материала по ID
func DBGetMaterialByID(materialID int) (*model.Material, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var material model.Material
	err := db.QueryRow("SELECT * FROM materials WHERE material_id = $1", materialID).Scan(
		&material.ID, &material.Lesson_id, &material.Title, &material.File_url, &material.TypeOfMaterial,
	)

	if err != nil {
		return nil, err
	}
	return &material, nil
}

// DBGetEnrollmentByID - получение записи по ID
func DBGetEnrollmentByID(enrollmentID int) (*model.Enrollment, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var enrollment model.Enrollment
	err := db.QueryRow("SELECT * FROM enrollments WHERE enrollment_id = $1", enrollmentID).Scan(
		&enrollment.ID, &enrollment.User_id, &enrollment.Course_id,
	)

	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}
func TakeTable(db *sql.DB, c *gin.Context, tableName string) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		rows.Scan(values...)
		for i, v := range values {
			if i > 0 {
				c.String(http.StatusOK, "|")
			}
			c.String(http.StatusOK, fmt.Sprintf("%v", *(v.(*interface{}))))
		}
		c.String(http.StatusOK, "\n")
	}
	return nil
}
func DBFindUserByGoogleID(googleID string) (*model.User, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var user model.User
	err := db.QueryRow("SELECT user_id, full_name, username, email, role, google_id, avatar FROM users WHERE google_id = $1", googleID).
		Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Role, &user.GoogleID, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DBFindUserByYandexID(yandexID string) (*model.User, error) {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var user model.User
	err := db.QueryRow("SELECT user_id, full_name, username, email, role, yandex_id, avatar FROM users WHERE yandex_id = $1", yandexID).
		Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Role, &user.YandexID, &user.Avatar)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// pkg/db_database/databasepostgres.go
func DBCreateUserFromOAuth(user *model.User) error {
	db := DBConnect(model.ConnStrUsers)
	defer db.Close()

	var userID int

	// Используем "student" как роль по умолчанию для OAuth пользователей
	role := "student"
	if user.Role != "" {
		// Проверяем, что роль соответствует допустимым значениям
		validRoles := map[string]bool{"student": true, "teacher": true, "admin": true}
		if validRoles[user.Role] {
			role = user.Role
		}
	}

	passwordHash := "oauth_user_no_password"
	if user.PasswordFirst != "" {
		var err error
		passwordHash, err = JSONJWT.HashPassword(user.PasswordFirst)
		if err != nil {
			return err
		}
	}

	err := db.QueryRow(`
		INSERT INTO users (full_name, email, google_id, yandex_id, avatar, role, username, password_hash) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING user_id`,
		user.FullName,
		user.Email,
		user.GoogleID,
		user.YandexID,
		user.Avatar,
		role, // Используем проверенную роль
		user.Username,
		passwordHash,
	).Scan(&userID)

	if err != nil {
		log.Printf("Error creating OAuth user: %v", err)
		return err
	}

	user.ID = userID
	user.Role = role // Обновляем роль в объекте пользователя
	return nil
}
