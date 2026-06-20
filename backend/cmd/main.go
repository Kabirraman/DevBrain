package main

import (
	"log"

	"github.com/Kabirraman/DevBrain/internal/auth"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/resources"
	"github.com/Kabirraman/DevBrain/internal/middleware"
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

	// Public Routes
	router.POST(
		"/api/auth/register",
		auth.RegisterHandler,
	)

	router.POST(
	"/api/auth/login",
	auth.LoginHandler,
)

	// Protected Routes
protected := router.Group("/api")

protected.Use(
	middleware.AuthMiddleware(),
)

protected.GET(
	"/me",
	auth.MeHandler,
)

protected.POST(
	"/resources",
	resources.CreateResourceHandler,
)

protected.GET(
	"/resources",
	resources.GetResourcesHandler,
)

protected.GET(
	"/resources/:id",
	resources.GetResourceHandler,
)

protected.POST(
	"/resources/blog",
	resources.ImportBlogHandler,
)
log.Println("Server running on :8080")

	router.Run(":8080")
}

