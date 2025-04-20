package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"database/sql"
	"project/database" // Импорт пакета database
	"project/models"   // Импорт пакета models
	"project/utils"    // Import the utils package
)

// CreateUserHandler обрабатывает запрос на создание нового пользователя.
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("createUserHandler called")
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			fmt.Println("Method not allowed")
			return
		}

		var user models.User // Используем models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Printf("JSON Decode error: %v\n", err)
			return
		}

		fmt.Printf("Received user %+v\n", user)

		hashedPassword, _ := utils.HashPassword(user.PasswordHash) // Use utils.HashPassword
		user.PasswordHash = hashedPassword

		database.AddData(db, "INSERT INTO users (full_name, Username, Email, Password_Hash, Role) VALUES ($1, $2, $3, $4, $5)", user)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Пользователь успешно создан")
		fmt.Println("User created successfully")
	}
}
