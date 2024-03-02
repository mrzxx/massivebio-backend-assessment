package main

import (
	"fmt"
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
	fmt.Println("Server listing 8080 port")
	http.ListenAndServe(":8080", nil)

}
