package database

import (
	"backend/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

func ConnectDb() {
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file")
	}
	// Environment değişkenlerini oku
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	//timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Printf("[error] failed to initialize database, got error %v\n", err)
		println("Could not connect to the database")
	} else {
		println("Connected to the database")
	}

	DB = db

	err = db.AutoMigrate(
		&models.Poem{},
		&models.Admin{},
		&models.Log{},
		&models.Book{},
		&models.Comment{},
		&models.Reminder{},
		&models.Homepage{},
		&models.MihrimahCard{},
		&models.Friendship{},
	)
	if err != nil {
		panic("Could not migrate to the database")
	} else {
		println("Migrated to the database")
	}
}
