package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anandaa033/thl-weather-service/handlers"
)

func main() {
	http.HandleFunc("/api/weather", handlers.WeatherHandler)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
