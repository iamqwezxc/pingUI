package main

import (
	"fmt"
	"log"
	"net/http"
	"project/database"
	"project/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := database.Connect("user=postgres password=psql dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/users", handlers.CreateUserHandler(db)).Methods("POST")

	fmt.Println("Сервер запущен на порту 8280")
	log.Fatal(http.ListenAndServe(":8280", router))

}
