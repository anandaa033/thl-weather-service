package server

import (
	"fmt"
	"log"
	"net/http"

	handlers "thlWeatherService/handlers/handler_weather"
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

	http.HandleFunc("/api/weather", handlers.WeatherHandler)
	http.HandleFunc("/api/weather/multiple", handlers.MultiCityWeatherHandler)
	http.HandleFunc("/api/user/cities", handlers.UserCitiesHandler)

	port := config.Get().App.Config.Port
	log.Printf("✅ Server is running at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
