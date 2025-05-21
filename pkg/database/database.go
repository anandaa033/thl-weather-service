package database

import (
	"fmt"
	"log"

	"thlWeatherService/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	conf := config.Get()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Gorm.Host,
		conf.Database.Gorm.Port,
		conf.Database.Gorm.Username,
		conf.Database.Gorm.Password,
		conf.Database.Gorm.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return err
	}

	DB = db
	log.Println("✅ Database connection established")
	return nil
}
