package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/Kabirraman/DevBrain/internal/database"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	err = database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected Successfully ")
}