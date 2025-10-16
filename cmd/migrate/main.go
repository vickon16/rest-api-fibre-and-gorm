package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/database"
	"github.com/vickon16/rest-api-fibre-and-gorm/cmd/models"
)

func main() {
	_ = godotenv.Load()
	// Connect to DB
	database.ConnectDb()
	db := database.Database.Db

	log.Println("🚀 Running migrations...")

	log.Println("⚠️ Dropping all tables...")

	// Drop tables
	// if err := db.Migrator().DropTable(
	// 	&models.User{},
	// 	&models.Product{},
	// 	&models.Order{},
	// ); err != nil {
	// 	log.Fatal("❌ Failed to drop tables:", err)
	// }

	// log.Println("✅ Tables dropped successfully.")
	// log.Println("🚀 Recreating tables...")

	// Add all models that should be migrated
	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
	); err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ Migration completed successfully.")
}

// go run cmd/migrate/main.go
