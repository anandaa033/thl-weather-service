package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/anandaa033/thl-weather-service/models"
)

const apiKey = "e5ae287b7d3c0a860ae06db0e2ff0f1d"

func FetchWeatherData(city string) (*models.ForecastResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric",
		city, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", body)
	}

	var forecast models.ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}
	return &forecast, nil
}

func CalculateAverageTemp(forecast *models.ForecastResponse) float64 {
	if len(forecast.List) == 0 {
		return 0
	}
	var total float64
	for _, item := range forecast.List {
		total += item.Main.Temp
	}
	return total / float64(len(forecast.List))
}

func SummarizeDailyForecasts(forecast *models.ForecastResponse) []map[string]interface{} {
	dailyData := make(map[string][]models.ForecastItem)
	for _, item := range forecast.List {
		date := item.DtTxt[:10]
		dailyData[date] = append(dailyData[date], item)
	}

	dates := make([]string, 0, len(dailyData))
	for date := range dailyData {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02", dates[i])
		t2, _ := time.Parse("2006-01-02", dates[j])
		return t1.Before(t2)
	})

	var results []map[string]interface{}
	for _, date := range dates {
		items := dailyData[date]
		minTemp := items[0].Main.Temp
		maxTemp := items[0].Main.Temp
		mainCondition := items[0].Weather[0].Main

		for _, entry := range items {
			if entry.Main.Temp < minTemp {
				minTemp = entry.Main.Temp
			}
			if entry.Main.Temp > maxTemp {
				maxTemp = entry.Main.Temp
			}
		}

		parsedDate, _ := time.Parse("2006-01-02", date)
		weekday := parsedDate.Weekday().String()

		results = append(results, map[string]interface{}{
			"day":       weekday,
			"condition": mainCondition,
			"temp_min":  minTemp,
			"temp_max":  maxTemp,
		})
	}

	return results
}
