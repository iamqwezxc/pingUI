package main

import (
	_ "github.com/lib/pq"

	model "github.com/iamqwezxc/pingUI/Backend/models"
	database "github.com/iamqwezxc/pingUI/Backend/pkg/db_database"

	//JWTs "studyproject/pkg/json_jwt"
	website "github.com/iamqwezxc/pingUI/Backend/pkg/wb_website"
)

func main() {

	//http.HandleFunc("/user", JWTs.AuthMiddleware(JWTs.ProtectedHandler))
	database.DBConnect(model.ConnStrUsers)
	website.WBStarsWebSite()

}
