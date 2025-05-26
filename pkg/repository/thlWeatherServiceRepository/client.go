package thlWeatherService

import (
	// "thlWeatherService/pkg/repository/thlWeatherServiceRepository/models"
	"thlWeatherService/pkg/utility/gormDB"
	"thlWeatherService/pkg/repository/thlWeatherServiceRepository/models"

	"gorm.io/gorm"
)

type Client struct {
	database *gorm.DB
}

func New(host string, port int, username, password, dbname string) (*Client, error) {

	var err error
	client, err := gormDB.New(host, port, username, password, dbname)
	if err != nil {
		return &Client{}, err
	}

	db := client.GetDB()

	// Automigrate multiple models
	err = db.AutoMigrate(
		models.User{},
	) // Add more models as needed
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		database: db,
	}, nil
}
