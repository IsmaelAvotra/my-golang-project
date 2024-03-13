package main

import (
	"log"
	"os"

	"github.com/IsmaelAvotra/pkg/api"
	"github.com/IsmaelAvotra/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file:", err)
		}
	}

	err := godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.ConnectDatabase()

	gin.SetMode(gin.DebugMode)

	r := api.InitRouter()

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
