package model

const ConnStrUsers = "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable"

type User struct {
	ID             int    `json:"ID"`
	FullName       string `json:"Full_Name"`
	Username       string `json:"Username"`
	Email          string `json:"Email"`
	PasswordFirst  string `json:"PasswordFHash"`
	PasswordSecond string `json:"PasswordSHash"`
	Role           string `json:"Role"`
}
