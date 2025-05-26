package handlers

import (
	"encoding/json"
	"net/http"

	"thlWeatherService/pkg/services"
)

func MultiCityWeatherHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Missing user_id",
		})
		return
	}

	cities := services.GetUserCities(userID)
	if len(cities) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "no cities",
			"results": []interface{}{},
		})
		return
	}

	var results []map[string]interface{}
	for _, city := range cities {
		forecast, err := services.FetchWeatherData(city)
		if err != nil {
			continue
		}

		avgTempC := services.CalculateAverageTemp(forecast)

		avgTempF := (avgTempC * 9.0 / 5.0) + 32.0

		conditionCount := make(map[string]int)
		for _, item := range forecast.List {
			if len(item.Weather) > 0 {
				condition := item.Weather[0].Main
				conditionCount[condition]++
			}
		}

		mostCommonCondition := ""
		maxCount := 0
		for cond, count := range conditionCount {
			if count > maxCount {
				mostCommonCondition = cond
				maxCount = count
			}
		}

		results = append(results, map[string]interface{}{
			"city":                city,
			"country":             forecast.City.Country,
			"avg_temp_celsius":    avgTempC,
			"avg_temp_fahrenheit": avgTempF,
			"condition":           mostCommonCondition,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "success",
		"all_city": results,
	})
}
