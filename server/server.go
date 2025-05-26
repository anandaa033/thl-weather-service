package server

import (
	"fmt"
	"log"
	"net/http"

	handlersLocation "thlWeatherService/server/handlers/handler_location"
	handlersWheather "thlWeatherService/server/handlers/handler_weather"
	"thlWeatherService/pkg/config"
	"thlWeatherService/pkg/database"
)

func Run() {
	if err := config.Init("/config/config/config.yaml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	http.HandleFunc("/api/weather", handlersWheather.WeatherHandler)
	http.HandleFunc("/api/weather/multiple", handlersLocation.MultiCityWeatherHandler)
	http.HandleFunc("/api/user/cities", handlersLocation.UserCitiesHandler)

	port := config.Get().App.Config.Port
	log.Printf("✅ Server is running at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
