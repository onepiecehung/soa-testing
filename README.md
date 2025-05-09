# Online Store Product Management System

## Overview
This is a RESTful API for managing products in an online store. The system provides endpoints for product management, including CRUD operations, filtering, and search functionality. It's built with performance and scalability in mind.

## Features Implemented in This Project
### Security
✅ JWT Authentication: Implemented comprehensive authentication with access and refresh tokens

✅ Security Threat Protection: 
- SQL Injection: Implemented using GORM with prepared statements
- XSS: Implemented XSSMiddleware
- CSRF: Implemented CSRFMiddleware
- DDOS: Implemented RateLimitMiddleware

✅ Data Encryption: Implemented password hashing using bcrypt

### Maintainability
✅ Code Structure: Organized following Clean Architecture principles with distinct layers:
- Models
- Repositories
- Services
- Handlers
- Middleware
- Types/DTOs

### Database Optimization
✅ Relationship Management: Properly implemented database relationships:
- Product-Category (Many-to-Many)
- Product-Review (One-to-Many)
- User-Wishlist (One-to-Many)
- User-Review (One-to-Many)

✅ Query Optimization: Leveraged GORM features for optimal performance:
- Eager loading
- Preload
- Pagination

### Unit Tests
❌ Unit Testing: Pending implementation of comprehensive unit tests for API endpoints

### Logging
✅ Comprehensive Logging System: Implemented using logrus with the following features:
- Request Logging: Basic request logging through middleware
- Error Logging: Basic error tracking in handlers
- Text Format: Logs formatted in text format
- Log Levels: Basic log level support (Info, Error)
- Middleware: Basic logging middleware

### Performance
✅ Pagination: Implemented across all API endpoints

✅ Query Optimization: Leveraged GORM's optimization features

✅ Database Level Caching: Implemented through:
- Prepared Statements
- Connection Pooling
- Query Optimization

❌ Distributed Caching: Pending implementation of Redis/Memcached caching strategy

### Scalability
✅ Architecture: Designed with scalability in mind:
- Separation of concerns
- Dependency injection
- Interface-based design

❌ Microservices: Pending implementation of microservices architecture

## Technologies Used
- Go (Golang) - Backend programming language
- PostgreSQL - Database
- GORM - ORM for database operations
- JWT - Authentication
- Docker - Containerization
- Gin - Web framework

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Generating Swagger Documentation](#generating-swagger-documentation)

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 15 or higher
- Swag CLI (for generating Swagger documentation)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/onepiecehung/soa-testing.git
cd soa-testing
```

2. Install dependencies:
```bash
go mod download
```

3. Install Swag CLI (for generating Swagger documentation):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Setup Instructions

### Environment Setup
1. Create a `.env` file in the root directory with the following variables:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=product_management
JWT_SECRET=your_jwt_secret
JWT_REFRESH_SECRET=your_jwt_refresh_secret
JWT_EXPIRATION=24h
PORT=8080
ENVIRONMENT=development
CSRF_SECRET=your_csrf_secret
RATE_LIMIT=100
RATE_WINDOW=1h
```

### Running the Application
1. Start the database:
```bash
docker-compose up -d db
```

2. Run migrations:
```bash
go run cmd/migrate/main.go
```

3. Start the application:
```bash
go run cmd/server/main.go
```

### Running with Docker
```bash
docker-compose up --build
```

## Account test
#### User
```
email: "user@soa.com"
PW: "password1213"
```
#### Admin
```
email: "admin@soa.com"
PW: "password1213"
```
<strong><span style="color: red;">Note: I initialized the list of test products and accounts when setting up or starting the project.</span></strong>

## API Documentation

The API documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Swagger JSON: `http://localhost:8080/swagger/doc.json`

![SWAGGER](./assets/images/swagger.png)

## Generating Swagger Documentation

### Initial Setup
1. Install Swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate Swagger documentation:
```bash
swag init -g cmd/server/main.go
```

This will create a `docs` directory with the generated Swagger documentation.

### Updating Documentation
When you make changes to your API endpoints or documentation comments, you need to regenerate the Swagger documentation:

```bash
swag init -g cmd/server/main.go
```

### Adding Documentation to Your Code

Add Swagger annotations to your code using comments. Here's an example:

```go
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product true "Product object"
// @Success 201 {object} Product
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/products [post]
func CreateProduct(c *gin.Context) {
    // Handler implementation
}
```

### Available Swagger Annotations

- `@Summary`: Brief description of the endpoint
- `@Description`: Detailed description
- `@Tags`: Group endpoints by tags
- `@Accept`: Accepted request content type
- `@Produce`: Response content type
- `@Param`: Request parameters
- `@Success`: Success response
- `@Failure`: Error responses
- `@Router`: Endpoint path and HTTP method

For more detailed information about Swagger annotations, refer to the [Swaggo documentation](https://github.com/swaggo/swag#declarative-comments-format).


## Database Schema

### Users
- id (SERIAL PRIMARY KEY)
- username (VARCHAR(255) UNIQUE)
- email (VARCHAR(255) UNIQUE)
- password_hash (VARCHAR(255))
- full_name (VARCHAR(255))
- role (VARCHAR(50))
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Products
- id (SERIAL PRIMARY KEY)
- name (VARCHAR(255))
- description (TEXT)
- price (DECIMAL(10,2))
- stock_quantity (INT)
- status (VARCHAR(50))
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Categories
- id (SERIAL PRIMARY KEY)
- name (VARCHAR(255))
- description (TEXT)

### Product_Categories
- product_id (INT)
- category_id (INT)
- PRIMARY KEY (product_id, category_id)

### Reviews
- id (SERIAL PRIMARY KEY)
- product_id (INT)
- user_id (INT)
- rating (INT)
- comment (TEXT)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Wishlist
- user_id (INT)
- product_id (INT)
- added_at (TIMESTAMP)
- PRIMARY KEY (user_id, product_id)

## Database Migration

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 15 or higher
- GORM installed in your project

### Step 1: Create Migration Files

1. Create a new directory for migrations:
```bash
mkdir -p migrations
```

2. Create your first migration file (e.g., `migrations/0001_initial_schema.sql`):
```sql
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INT NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create product_categories table
CREATE TABLE IF NOT EXISTS product_categories (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

-- Create reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create wishlist table
CREATE TABLE IF NOT EXISTS wishlist (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, product_id)
);
```

### Step 2: Create Migration Script

1. Create a new file `cmd/migrate/main.go`:
```go
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"product-management/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get list of migration files
	migrationFiles, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	// Sort migration files by name
	sort.Strings(migrationFiles)

	// Apply migrations
	for _, file := range migrationFiles {
		// Check if migration has already been applied
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", filepath.Base(file)).Scan(&count)
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping...", file)
			continue
		}

		// Read migration file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		// Split file into individual statements
		statements := strings.Split(string(content), ";")
		
		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to start transaction: %v", err)
		}

		// Execute each statement
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to execute statement in %s: %v\nStatement: %s", file, err, stmt)
			}
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", filepath.Base(file))
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to record migration: %v", err)
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		log.Printf("Successfully applied migration: %s", file)
	}

	log.Println("All migrations completed successfully")
}
```

### Step 3: Running Migrations

1. Make sure your database is running:
```bash
docker compose up -d db
```

2. Run the migration script:
```bash
go run cmd/migrate/main.go
```

### Step 4: Creating New Migrations

1. Create a new migration file with a sequential number:
```bash
touch migrations/0002_add_new_feature.sql
```

2. Add your SQL statements to the new file:
```sql
-- Add new column to products table
ALTER TABLE products ADD COLUMN new_column VARCHAR(255);

-- Create new table
CREATE TABLE new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
```

3. Run the migration script again:
```bash
go run cmd/migrate/main.go
```

### Step 5: Rolling Back Migrations

To roll back a specific migration:

1. Create a rollback file with the same number but with "_rollback" suffix:
```bash
touch migrations/0002_add_new_feature_rollback.sql
```

2. Add rollback SQL statements:
```sql
-- Rollback changes
ALTER TABLE products DROP COLUMN new_column;
DROP TABLE new_table;
```

3. Run the rollback:
```bash
go run cmd/migrate/main.go --rollback 0002_add_new_feature
```

### Best Practices

1. Always create migration files with sequential numbers
2. Include both up and down migrations
3. Test migrations in a development environment first
4. Backup your database before running migrations in production
5. Use transactions for atomic migrations
6. Document your migrations in the README
7. Keep migration files in version control

### Troubleshooting

1. If a migration fails:
   - Check the error message
   - Fix the SQL statements
   - Run the migration again

2. If you need to reset the database:
   - Drop all tables
   - Delete all records from the migrations table
   - Run migrations from scratch

3. Common issues:
   - Syntax errors in SQL
   - Missing dependencies
   - Permission issues
   - Database connection problems

## Automatic Migration Generation

### Using Golang-Migrate

Golang-Migrate is a database migration tool that can automatically generate migration files from your Go models.

#### Installation

1. Install the migrate CLI tool:
```bash
# For Windows
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# For Linux/Mac
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate.linux-amd64 /usr/local/bin/migrate
```

#### Setting Up Automatic Migrations

1. Create a new directory for your migrations:
```bash
mkdir -p migrations
```

2. Create a new file `cmd/migrate/main.go`:
```go
package main

import (
	"log"
	"os"
	"product-management/config"
	"product-management/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Category{},
		&models.Review{},
		&models.Wishlist{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("Auto migration completed successfully")
}
```

3. Create your models in `internal/models/models.go`:
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	FullName     string
	Role         string `gorm:"default:user"`
	Reviews      []Review    // One user can have many reviews
	Wishlists    []Wishlist  // One user can have many wishlist items
}

// Product represents a product in the store
type Product struct {
	gorm.Model
	Name          string  `gorm:"not null"`
	Description   string
	Price         float64 `gorm:"not null"`
	StockQuantity int     `gorm:"not null;default:0"`
	Status        string  `gorm:"default:active"`
	Reviews       []Review // One product can have many reviews
	Categories    []Category `gorm:"many2many:product_categories;"` // Many-to-Many relationship with categories
	Wishlists     []Wishlist // One product can be in many wishlists
}

// Category represents a product category
type Category struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Products    []Product `gorm:"many2many:product_categories;"` // Many-to-Many relationship with products
}

// Review represents a product review
type Review struct {
	gorm.Model
	ProductID uint    // Foreign key to Product
	Product   Product `gorm:"foreignKey:ProductID"` // Belongs to one Product
	UserID    uint    // Foreign key to User
	User      User    `gorm:"foreignKey:UserID"`    // Belongs to one User
	Rating    int     `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string
}

// Wishlist represents a user's wishlist item
type Wishlist struct {
	gorm.Model
	UserID    uint    // Foreign key to User
	User      User    `gorm:"foreignKey:UserID"`    // Belongs to one User
	ProductID uint    // Foreign key to Product
	Product   Product `gorm:"foreignKey:ProductID"` // Belongs to one Product
	AddedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ProductID  uint
	CategoryID uint
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}
```

4. Generate migration files:
```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq create_initial_tables

# This will create two files:
# - migrations/000001_create_initial_tables.up.sql
# - migrations/000001_create_initial_tables.down.sql
```

The generated SQL for the up migration will look like this:
```sql
-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock_quantity INT NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create reviews table
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create product_categories table (junction table for many-to-many relationship)
CREATE TABLE product_categories (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, category_id)
);

-- Create wishlist table
CREATE TABLE wishlists (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_id, product_id)
);
```

5. Run the auto migration:
```bash
go run cmd/migrate/main.go
```

#### Benefits of Using Golang-Migrate

1. **Version Control**: Each migration is versioned and tracked
2. **Up/Down Migrations**: Supports both applying and rolling back migrations
3. **Multiple Databases**: Supports PostgreSQL, MySQL, SQLite, and more
4. **CLI Tool**: Easy to use command-line interface
5. **Go Integration**: Can be used directly from Go code

#### Common Commands

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq migration_name

# Apply all pending migrations
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up

# Rollback the last migration
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" version
```

#### Best Practices for Automatic Migrations

1. Always review the generated SQL before applying migrations
2. Keep your models in sync with your database schema
3. Use transactions for complex migrations
4. Test migrations in a development environment first
5. Backup your database before running migrations in production
6. Document schema changes in your code comments

## Building Docker Container
```bash
docker build -t hungnh-product-management .
```

## Security Features
- JWT Authentication with token expiration
- Password hashing using bcrypt
- SQL injection prevention through GORM
- XSS protection middleware
- CSRF protection
- Rate limiting
- Input validation
- Secure cookie settings
- HTTPS enforcement

## Monitoring and Logging
The application includes built-in logging using logrus. Logs are written to stdout and can be collected by your logging infrastructure.

### Log Levels
- DEBUG: Detailed information for debugging
- INFO: General operational information
- WARN: Warning messages for potential issues
- ERROR: Error messages for failed operations
- FATAL: Critical errors that require immediate attention

### Log Format
```json
{
  "level": "info",
  "time": "2024-01-01T12:00:00Z",
  "message": "Server started",
  "port": 8080
}
```

## Performance Optimization
- Database query optimization using GORM
- Pagination for large datasets
- Connection pooling for database connections
- Rate limiting to prevent abuse

## Scaling Considerations
- Horizontal scaling with multiple instances
- Load balancing using Nginx or similar
- Database replication for read operations
- Caching layer for frequently accessed data
- Message queue for async operations

## Troubleshooting
Common issues and solutions:

1. Database Connection Issues
```bash
# Check database connection
psql -h localhost -U postgres -d product_management

# Verify database credentials
cat .env
```

2. JWT Authentication Issues
```bash
# Verify JWT secret
echo $JWT_SECRET

# Check token expiration
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/products
```

3. Rate Limiting Issues
```bash
# Check rate limit settings
cat .env | grep RATE

# Monitor request count
curl -H "X-Request-ID: $(uuidgen)" http://localhost:8080/api/v1/products
```
