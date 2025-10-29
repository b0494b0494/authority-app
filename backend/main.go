package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"authority-app/backend/internal/models"
	"authority-app/backend/internal/openfga"
	"authority-app/backend/internal/repository"
	"authority-app/backend/internal/service"
	"authority-app/backend/internal/util"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

		token, err := util.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}

func main() {
	// Initialize database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	// Seed a default user if none exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, err := util.HashPassword("password")
		if err != nil {
			log.Fatalf("failed to hash default password: %v", err)
		}
		defaultUser := models.User{Username: "testuser", Password: hashedPassword}
		if err := db.Create(&defaultUser).Error; err != nil {
			log.Fatalf("failed to create default user: %v", err)
		}
		log.Println("Default user 'testuser' created with password 'password'")
	}

	// Initialize OpenFGA client
	fgaClient, err := openfga.NewFGAClient()
	if err != nil {
		log.Fatalf("failed to initialize OpenFGA client: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	uploadService := service.NewUploadService(fgaClient)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User registration endpoint
	r.POST("/register", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := authService.Register(input.Username, input.Password); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User registered successfully"})
	})

	// User login endpoint
	r.POST("/login", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		tokenString, err := authService.Login(input.Username, input.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	})

	// Protected endpoint
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(200, gin.H{"message": "Welcome to the protected area, " + username + "!"})
	})

	// Upload endpoint
	r.POST("/upload", AuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to get file"})
			return
		}

		if err := uploadService.UploadAndExtract(file, username); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "File uploaded and extracted successfully! Owner tuple created."})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}