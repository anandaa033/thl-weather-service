package services

import (
	"errors"
	"log"
	"time"

	"thlWeatherService/pkg/database"
	"thlWeatherService/pkg/repository/thlWeatherServiceRepository/models"

	"gorm.io/gorm"
)

func AddCityToUser(userID, city string) {
	db := database.DB

	var existing models.UserCity
	err := db.Unscoped().
		Where("user_id = ? AND city = ?", userID, city).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var maxOrder int
			db.Model(&models.UserCity{}).
				Where("user_id = ? AND deleted_at IS NULL", userID).
				Select("COALESCE(MAX(city_order), -1)").Scan(&maxOrder)

			entry := models.UserCity{
				UserID:    userID,
				City:      city,
				CityOrder: maxOrder + 1,
			}
			if err := db.Create(&entry).Error; err != nil {
				log.Println("AddCityToUser → create error:", err)
			}
		} else {
			log.Println("AddCityToUser → db error:", err)
		}
		return
	}

	if existing.DeletedAt != nil {
		var maxOrder int
		db.Model(&models.UserCity{}).
			Where("user_id = ? AND deleted_at IS NULL", userID).
			Select("COALESCE(MAX(city_order), -1)").Scan(&maxOrder)

		existing.DeletedAt = nil
		existing.UpdatedAt = time.Now()
		existing.CityOrder = maxOrder + 1
		if err := db.Save(&existing).Error; err != nil {
			log.Println("AddCityToUser → restore error:", err)
		}
	}
}

func RemoveCityFromUser(userID, city string) {
	db := database.DB

	if err := db.Model(&models.UserCity{}).
		Where("user_id = ? AND city = ? AND deleted_at IS NULL", userID, city).
		Update("deleted_at", time.Now()).Error; err != nil {
		log.Println("RemoveCityFromUser error:", err)
	}
}

func GetUserCities(userID string) []string {
	db := database.DB
	var results []models.UserCity

	err := db.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("city_order asc").
		Find(&results).Error
	if err != nil {
		log.Println("GetUserCities error:", err)
		return nil
	}

	var cities []string
	for _, row := range results {
		cities = append(cities, row.City)
	}
	return cities
}

func ReorderCities(userID string, order []string) {
	db := database.DB
	now := time.Now()

	for i, city := range order {
		err := db.Model(&models.UserCity{}).
			Where("user_id = ? AND city = ? AND deleted_at IS NULL", userID, city).
			Updates(map[string]interface{}{
				"city_order": i,
				"updated_at": now,
			}).Error
		log.Printf("ReorderCities: updating %s to position %d", city, i)
		if err != nil {
			log.Println("ReorderCities error:", err)
		}
	}
}
