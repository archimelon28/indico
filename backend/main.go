package main

import (
	database "backend/database"
	"backend/models"
	"backend/routes"
	"log"
	"net/http"
)

func main() {
	dbClient := database.ConnectDatabase()

	databaseClient := models.NewDBClient(dbClient)

	r := routes.SetupRouter(databaseClient)

	log.Fatal(http.ListenAndServe(":8080", r))
}
