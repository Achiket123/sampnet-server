package database

import (
	"fmt"
	"log"
	"os"
	"server/internal/platform/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	DB = initDB()
}

func initDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatalf("Database configuration missing. Please check environment variables.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func PerformMigrations() error {
	if err := DB.AutoMigrate(&models.Organisation{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.UserModel{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Task{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.File{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Employee{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Attendance{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Project{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Team{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.TeamMember{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Notification{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.Chat{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.ChatMessage{}); err != nil {
		return err
	}
	if err := DB.AutoMigrate(&models.CallState{}); err != nil {
		return err
	}
	return nil
}
