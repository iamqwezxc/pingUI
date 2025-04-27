package main

import (
	_ "github.com/lib/pq"

	model "studyproject/models"
	database "studyproject/pkg/db_database"

	//JWTs "studyproject/pkg/json_jwt"
	website "studyproject/pkg/wb_website"
)

func main() {

	//http.HandleFunc("/user", JWTs.AuthMiddleware(JWTs.ProtectedHandler))
	//http.ListenAndServe(":8080", nil)

	database.DBConnect(model.ConnStrUsers)
	website.WBStarsWebSite()

}
