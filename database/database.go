package database

import (
	"database/sql"
	"fmt"
	"log"
	"project/models"

	_ "github.com/lib/pq"
)

func Connect(srokaSQLconn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", srokaSQLconn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Успешно подключено к PostgreSQL!")
	return db, nil
}

func AddData(db *sql.DB, str string, u models.User) {
	stmt, err := db.Prepare(str)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.FullName, u.Username, u.Email, u.PasswordHash, u.Role)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные успешно вставлены!")
}
