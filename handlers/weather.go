package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/anandaa033/thl-weather-service/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "Missing 'city' parameter", http.StatusBadRequest)
		return
	}

	forecast, err := services.FetchWeatherData(city)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching weather data: %v", err), http.StatusInternalServerError)
		return
	}

	avgTemp := services.CalculateAverageTemp(forecast)
	today := time.Now().Format("2006-01-02")
	var todayForecast []map[string]interface{}
	var sumHumidity, sumWind, sumPrecip float64
	var count int

	var (
		totalDayTemp, totalNightTemp float64
		dayCount, nightCount         int
		maxDayTemp, maxNightTemp     float64
		minDayTemp, minNightTemp     float64
	)

	const maxExpectedRainMM = 30.0

	for _, item := range forecast.List {
		dateStr := item.DtTxt[:10]
		if dateStr == today {
			sumHumidity += float64(item.Main.Humidity)
			sumWind += item.Wind.Speed

			if item.Rain != nil {
				sumPrecip += item.Rain.ThreeH
			}

			parsedTime, _ := time.Parse("2006-01-02 15:04:05", item.DtTxt)
			isDaytime := parsedTime.Hour() >= 6 && parsedTime.Hour() < 18

			if isDaytime {
				totalDayTemp += item.Main.Temp
				dayCount++
				if item.Main.Temp > maxDayTemp {
					maxDayTemp = item.Main.Temp
				}
				if item.Main.Temp < minDayTemp || minDayTemp == 0 {
					minDayTemp = item.Main.Temp
				}
			} else {
				totalNightTemp += item.Main.Temp
				nightCount++
				if item.Main.Temp > maxNightTemp {
					maxNightTemp = item.Main.Temp
				}
				if item.Main.Temp < minNightTemp || minNightTemp == 0 {
					minNightTemp = item.Main.Temp
				}
			}

			todayForecast = append(todayForecast, map[string]interface{}{
				"time":      item.DtTxt,
				"temp":      item.Main.Temp,
				"condition": item.Weather[0].Main,
			})
			count++
		}
	}

	var (
		avgDayTemp, avgNightTemp, humidityPct, windSpeedKmh, precipPct float64
		todayCondition                                                 string
	)

	if dayCount > 0 {
		avgDayTemp = totalDayTemp / float64(dayCount)
	}
	if nightCount > 0 {
		avgNightTemp = totalNightTemp / float64(nightCount)
	}

	if count > 0 {
		humidityPct = sumHumidity / float64(count)
		windSpeedKmh = sumWind / float64(count) * 3.6
		precipPct = (sumPrecip / maxExpectedRainMM) * 100
		if precipPct > 100 {
			precipPct = 100
		}
		precipPct = math.Round(precipPct*10) / 10
		todayCondition = todayForecast[0]["condition"].(string)
	}

	dailyForecasts := services.SummarizeDailyForecasts(forecast)

	response := map[string]interface{}{
		"city":            forecast.City.Name,
		"country":         forecast.City.Country,
		"avg_temp":        avgTemp,
		"humidity_pct":    humidityPct,
		"wind_speed_kmh":  windSpeedKmh,
		"precip_pct":      precipPct,
		"condition":       todayCondition,
		"today_forecasts": todayForecast,
		"daily_forecasts": dailyForecasts,
		"day": []map[string]interface{}{
			{
				"avg_day_temp": avgDayTemp,
				"max_day_temp": maxDayTemp,
				"min_day_temp": minDayTemp,
			},
		},
		"night": []map[string]interface{}{
			{
				"avg_night_temp": avgNightTemp,
				"max_night_temp": maxNightTemp,
				"min_night_temp": minNightTemp,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
