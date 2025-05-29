package pkg

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/iamqwezxc/pingUI/Backend/models"
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
