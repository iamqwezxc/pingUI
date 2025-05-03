package databasepostgres

import (
	"database/sql"
	"fmt"
	"log"

	model "github.com/iamqwezxc/pingUI/Backend/models"
)

func DBTakeTable() {
	db := DBConnect(model.ConnStrUsers)
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
		fmt.Println(user_id, full_name, Username, Email, Password_Hash, Role)
	}
	defer db.Close()

}

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
		user.Password,
		user.Role,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
