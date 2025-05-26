package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"thlWeatherService/pkg/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		writeError(w, http.StatusBadRequest, "Missing 'city' parameter")
		return
	}

	forecast, err := services.FetchWeatherData(city)
	if err != nil {
		if strings.Contains(err.Error(), "city not found") {
			writeError(w, http.StatusNotFound, "City not found")
		} else {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching weather data: %v", err))
		}
		return
	}

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")

	var todayForecast []map[string]interface{}
	var sumHumidity, sumWind, sumPrecip float64
	var count int
	var TempNow float64

	var (
		totalDayTemp, totalNightTemp float64
		dayCount, nightCount         int
		maxDayTemp, maxNightTemp     float64
		minDayTemp, minNightTemp     = math.MaxFloat64, math.MaxFloat64
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
				if item.Main.Temp < minDayTemp {
					minDayTemp = item.Main.Temp
				}
			} else {
				totalNightTemp += item.Main.Temp
				nightCount++
				if item.Main.Temp > maxNightTemp {
					maxNightTemp = item.Main.Temp
				}
				if item.Main.Temp < minNightTemp {
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

		var closestForecast map[string]interface{}
		var latestPastForecastTime time.Time
		var earliestFutureForecast map[string]interface{}
		var earliestFutureTime time.Time

		for _, forecast := range todayForecast {
			forecastTimeStr := forecast["time"].(string)
			forecastTime, err := time.ParseInLocation("2006-01-02 15:04:05", forecastTimeStr, loc)
			if err != nil {
				continue
			}

			if forecastTime.Before(now) || forecastTime.Equal(now) {
				if forecastTime.After(latestPastForecastTime) {
					latestPastForecastTime = forecastTime
					closestForecast = forecast
				}
			} else {
				if earliestFutureTime.IsZero() || forecastTime.Before(earliestFutureTime) {
					earliestFutureTime = forecastTime
					earliestFutureForecast = forecast
				}
			}
		}

		if closestForecast == nil && earliestFutureForecast != nil {
			closestForecast = earliestFutureForecast
		}

		if closestForecast != nil {
			todayCondition = closestForecast["condition"].(string)
			TempNow = closestForecast["temp"].(float64)
		}
	}

	if minDayTemp == math.MaxFloat64 {
		minDayTemp = 0
	}
	if minNightTemp == math.MaxFloat64 {
		minNightTemp = 0
	}

	dailyForecasts := services.SummarizeDailyForecasts(forecast)

	response := map[string]interface{}{
		"message":         "success",
		"city":            forecast.City.Name,
		"country":         forecast.City.Country,
		"temp_now":        TempNow,
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

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
