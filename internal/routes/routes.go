package routes

import (
	"product-management/internal/handlers"
	"product-management/internal/middleware"
	"product-management/internal/models"
	"product-management/internal/repositories"
	"product-management/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title           Product Management API
// @version         1.0
// @description     A product management service with categories, reviews, and more.
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

// SetupRoutes configures all the routes for the application
func SetupRoutes(db *gorm.DB, r *gin.Engine) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	categoryService := services.NewCategoryService()
	reviewService := services.NewReviewService(reviewRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	authService := services.NewAuthService()
	authHandler := handlers.NewAuthHandler(userRepo, authService)

	// API version group
	api := r.Group("/api/v1")

	// Product routes
	products := api.Group("/products")
	products.Use(middleware.AuthMiddleware())
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
		products.GET("", productHandler.ListProducts)

		// Wishlist routes
		wishlist := products.Group("/wishlist")
		{
			wishlist.GET("", productHandler.GetWishlist)
			wishlist.POST("/:product_id", productHandler.AddToWishlist)
			wishlist.DELETE("/:product_id", productHandler.RemoveFromWishlist)
			wishlist.GET("/count", productHandler.GetTotalWishlistCount)
		}
	}

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		auth.PUT("/me", middleware.AuthMiddleware(), authHandler.UpdateUser)
		auth.PUT("/password", middleware.AuthMiddleware(), authHandler.UpdatePassword)
		auth.GET("/users/:id", middleware.AuthMiddleware(), authHandler.GetUserByID)
		auth.GET("/users", middleware.AuthMiddleware(), authHandler.ListUsers)
		auth.PUT("/users/:id/role", middleware.AuthMiddleware(), middleware.RequireRole(string(models.RoleAdmin)), authHandler.UpdateUserRole)
		auth.DELETE("/users/:id", middleware.AuthMiddleware(), middleware.RequireRole(string(models.RoleAdmin)), authHandler.DeleteUser)
	}

	// Review routes
	reviews := api.Group("/reviews")
	reviews.Use(middleware.AuthMiddleware())
	{
		reviews.POST("/", reviewHandler.CreateReview)
		reviews.GET("/", reviewHandler.SearchReviews)
		reviews.GET("/count", reviewHandler.GetTotalReviews)
		reviews.GET("/:id", reviewHandler.GetReviewByID)
		// reviews.GET("/product/:productId", reviewHandler.GetReviewsByProductID)
		// reviews.GET("/user/:userId", reviewHandler.GetReviewsByUserID)
		// reviews.PUT("/:id", reviewHandler.UpdateReview)
		reviews.DELETE("/:id", reviewHandler.DeleteReview)
		// reviews.GET("/product/:productId/rating", reviewHandler.GetProductRating)
		// reviews.GET("/product/:productId/count", reviewHandler.GetProductReviewCount)
	}

	// Category routes
	categories := api.Group("/categories")
	categories.Use(middleware.AuthMiddleware())
	{
		categories.POST("", categoryHandler.CreateCategory)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
		categories.GET("", categoryHandler.GetAllCategories)
		categories.GET("/distribution", categoryHandler.GetCategoryDistribution)

		// Category-Product relationship routes
		categoryProducts := categories.Group("/:id/products")
		{
			categoryProducts.GET("", categoryHandler.GetProductsByCategoryID)
			categoryProducts.POST("/:productId", categoryHandler.AddProductToCategory)
			categoryProducts.DELETE("/:productId", categoryHandler.RemoveProductFromCategory)
		}
	}
}
