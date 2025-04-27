package databasepostgres

import (
	"database/sql"
	"log"
	model "studyproject/models"
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
		user.Password,
		user.Role,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
