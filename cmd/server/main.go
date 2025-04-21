package main

import (
	"log"
	"product-management/config"
	"product-management/docs"
	"product-management/internal/middleware"
	"product-management/internal/routes"
	"product-management/pkg/database"
	"product-management/pkg/seeder"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Product Management API
// @version         1.0
// @description     A RESTful API for managing products in an online store.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Seed products initial data
	if err := seeder.SeedProducts(database.DB); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}
	// Seed users initial data
	if err := seeder.SeedUsers(database.DB); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Swagger documentation
	docs.SwaggerInfo.Title = "Product Management API"
	docs.SwaggerInfo.Description = "A RESTful API for managing products in an online store"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.AutoLogger())
	router.Use(middleware.ErrorHandlerMiddleware())
	// temporary comment auth middleware
	// router.Use(middleware.AuthMiddleware())

	// Setup all routes
	routes.SetupRoutes(database.DB, router)

	// Start server
	log.Printf("Server starting on port 8080...")
	log.Printf("Swagger documentation available at http://localhost:8080/swagger/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
