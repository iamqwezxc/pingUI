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

type Course struct {
	ID            int    `json:"ID"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Thumbnail_url string `json:"thumbnail_url"`
	Instructor_id int    `json:"instructor_id"`
}

type Lesson struct {
	ID           int    `json:"lesson_id"`
	Course_id    int    `json:"course_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Video_url    string `json:"video_url"`
	Lesson_order int    `json:"lesson_order"`
}
