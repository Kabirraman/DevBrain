package main

import (
	"log"

	"github.com/Kabirraman/DevBrain/internal/auth"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/resources"
	"github.com/Kabirraman/DevBrain/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/Kabirraman/DevBrain/internal/concepts"
	"github.com/Kabirraman/DevBrain/internal/graph"
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

	router.Use(cors.New(cors.Config{
	AllowOrigins: []string{
		"http://localhost:3000",
	},
	AllowMethods: []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
	},
	AllowHeaders: []string{
		"Origin",
		"Content-Type",
		"Authorization",
	},
}))

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

protected.POST(
	"/concepts/extract",
	concepts.ExtractConceptsHandler,
)

protected.POST(
	"/relationships/extract",
	concepts.ExtractRelationshipsHandler,
)

protected.GET(
	"/graph",
	graph.GetGraphHandler,
)
log.Println("Server running on :8080")

	router.Run(":8080")
}


