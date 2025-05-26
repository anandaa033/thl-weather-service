package repository

import (
	"context"
	"thlWeatherService/pkg/config"

	// "template/pkg/repository/game"
	"thlWeatherService/pkg/repository/thlWeatherServiceRepository"
)

type Repository struct {
	db *thlWeatherService.Client
}

func New(ctx context.Context, config *config.Config) (*Repository, error) {
	db1Client, err := thlWeatherService.New(config.Database.Gorm.Host, config.Database.Gorm.Port, config.Database.Gorm.Username, config.Database.Gorm.Password, config.Database.Gorm.Database)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db1Client,
		//db2: db2Client,
	}, nil
}

func (r *Repository) Db() *thlWeatherService.Client {
	return r.db
}
