package model

const ConnStrUsers = "host=localhost port=5432 user=postgres password=psql dbname=postgres sslmode=disable"

type User struct {
	ID       int    `json:"ID"`
	FullName string `json:"Full_Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"PasswordHash"`
	Role     string `json:"Role"`
}
