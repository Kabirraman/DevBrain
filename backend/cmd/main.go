package main

import (
	"log"

	"github.com/Kabirraman/DevBrain/internal/auth"
	"github.com/Kabirraman/DevBrain/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	err = database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST(
		"/api/auth/register",
		auth.RegisterHandler,
	)

	log.Println("Server running on :8080")

	router.Run(":8080")
}