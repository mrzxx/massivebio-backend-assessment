package main

import (
	"massivebio/controllers"
	"massivebio/database"
	"net/http"
)

func main() {
	//Init Database
	database.ConnectDB()
	defer database.DB.Close()

	//Routes
	http.HandleFunc("/assignment/query", controllers.MassiveFilter)

	//Run Server
	http.ListenAndServe(":8080", nil)
}
