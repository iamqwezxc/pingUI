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
	_, err := db.Exec(
		"INSERT INTO users (full_name, Username, Email, Password_Hash, Role) VALUES ($1, $2, $3, $4, $5)",
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordFirst,
		user.Role,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
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

func GetSlice(db *sql.DB, tableName string) ([]model.User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.PasswordFirst, &user.PasswordSecond, &user.Role)
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
