package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	// "gorm.io/gorm/schema"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance


func configureConnectionPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error retrieving underlying SQL DB: %v", err)
		return
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)           // Number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Recycle connections every hour

	log.Println("✅ Database connection pool configured successfully")
}

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			// Disable table name pluralization (optional)
			SingularTable: false,
			// Prevent converting to snake_case
			NoLowerCase: false,
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to database \n", err.Error())
	}

	// Configure the connection pool here
	configureConnectionPool(db)

	Database = DbInstance{Db: db}
	log.Println("Database connection successfully opened.")

}
