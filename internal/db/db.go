package db

import (
	"banking-api/internal/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		log.Fatal("DB_DSN not loaded. Ensure .env is loaded")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("[Error] failed to intialize database, got error %v", err)
	}

	// Auto-Migrate models
	err = db.AutoMigrate(&models.Customer{}, &models.Branch{}, &models.Account{}, &models.Transaction{}, &models.LoanPayment{}, &models.Beneficiary{})
	if err != nil {
		log.Fatalf("Automigrated failed: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}
