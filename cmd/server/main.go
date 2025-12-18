package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"banking-api/internal/db"
	"banking-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory:", err)
	}

	envPath := filepath.Join(dir, ".env")
	log.Printf("Loading .env from: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Could not load .env from %s: %v", envPath, err)
	} else {
		log.Println("Successfully loaded .env")
	}

	dsn := os.Getenv("DB_DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	if dsn == "" {
		log.Fatal("DB_DSN not loaded!")
	}

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not loaded!")
	}

	if port == "" {
		port = "8080"
	}

	dbConn := db.Connect()
	handlers.InitHandlers(dbConn)
	r := gin.Default()
	r.Use(gin.Logger()) // Logging middleware

	// CORS Middleware - Applied globally for all routes
	r.Use(corsMiddleWare())

	// public routes
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)

	// protected routes
	protected := r.Group("")
	protected.Use(authMiddleWare())

	// accounts
	protected.POST("/accounts", handlers.CreateAccount)
	protected.GET("/accounts", handlers.ListAccounts)
	protected.POST("/transfers/:from_id", handlers.Transfer)
	protected.POST("/deposits/:account_id", handlers.Deposit)
	protected.GET("/accounts/:id/statements", handlers.GetStatements)

	// Loans
	protected.POST("/loans", handlers.CreateLoan)
	protected.GET("/loans", handlers.ListLoans)
	protected.POST("/loans/:id/repay", handlers.MakePayment)
	protected.POST("/loans/:id/payments", handlers.ListPayments)

	// beneficaries
	protected.POST("/beneficiaries", handlers.AddBeneficiary)

	log.Printf("Server starting on %s", port)
	r.Run(":" + port)

}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Validate token via GetUserID in handlers
		c.Next()
	}
}

func corsMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// handle preflight options requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
