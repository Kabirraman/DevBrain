package main

import (
	"log"
	"os"

	"github.com/Kabirraman/DevBrain/internal/analytics"
	"github.com/Kabirraman/DevBrain/internal/auth"
	"github.com/Kabirraman/DevBrain/internal/chat"
	"github.com/Kabirraman/DevBrain/internal/concepts"
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/gaps"
	"github.com/Kabirraman/DevBrain/internal/graph"
	"github.com/Kabirraman/DevBrain/internal/middleware"
	"github.com/Kabirraman/DevBrain/internal/resources"
	"github.com/Kabirraman/DevBrain/internal/search"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("no .env file found, reading from environment")
	}

	err = database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	frontendURL := os.Getenv("FRONTEND_URL")

	allowedOrigins := []string{"http://localhost:3000"}

	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
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

	router.GET(
		"/api/concepts/:name",
		concepts.GetConceptDetailsHandler,
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
	protected.POST(
    "/chat",
    chat.ChatHandler,
)

	protected.GET(
		"/gaps",
		gaps.GetGapsHandler,
	)

	protected.GET(
		"/gaps/domains",
		gaps.GetDomainsHandler,
	)

	protected.POST(
		"/search",
		search.SearchHandler,
	)

	protected.GET(
		"/analytics",
		analytics.GetAnalyticsHandler,
	)

	log.Println("Server running on :8080")

	router.Run(":8080")
}
