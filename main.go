package main

import (
	"fmt"
	database "go-multitenancy-boilerplate/database"
	routers "go-multitenancy-boilerplate/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start database services and load master database.
	database.StartDatabaseServices()

	r := routers.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	type Job interface {
		Run()
	}

	// Starting the router instance
	if err := r.Run(":" + port); err != nil {
		fmt.Print(err)
	}
}
