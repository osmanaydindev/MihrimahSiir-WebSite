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

	// Many2many join tablolarını açık modellere bağla. AutoMigrate'ten ÖNCE
	// çağrılmak zorunda: sonrasında GORM ilişkiyi kendi ürettiği şemayla
	// kurar ve CreatedAt kolonu hiç eklenmez. Tablo adları değişmiyor,
	// mevcut veriye dokunulmuyor — sadece kolon ekleniyor.
	joins := []struct {
		model interface{}
		field string
		join  interface{}
	}{
		{&models.Admin{}, "AdminLikedPoems", &models.AdminLikedPoem{}},
		{&models.Admin{}, "AdminBookmarkPoems", &models.AdminBookmarkPoem{}},
		{&models.Admin{}, "UserBooksRead", &models.UserBookRead{}},
	}
	for _, j := range joins {
		if err := db.SetupJoinTable(j.model, j.field, j.join); err != nil {
			panic(fmt.Sprintf("Could not set up join table for %s: %v", j.field, err))
		}
	}

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
		&models.BookVisibility{},
		&models.BookRequest{},
		&models.MailLog{},
		&models.CommentLike{},
		&models.CommentSave{},
	)
	if err != nil {
		panic("Could not migrate to the database")
	} else {
		println("Migrated to the database")
	}

	// GORM tag'leriyle ifade edilemeyen indeksler
	EnsureIndexes()
}
