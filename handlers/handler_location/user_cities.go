package handlers

import (
	"encoding/json"
	"net/http"

	"thlWeatherService/services"
)

func UserCitiesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Missing user_id",
		})
		return
	}

	switch r.Method {
	case "POST":
		var payload struct {
			City string `json:"city"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid JSON",
			})
			return
		}

		if payload.City == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Missing city",
			})
			return
		}

		services.AddCityToUser(userID, payload.City)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
		})

	case "DELETE":
		city := r.URL.Query().Get("city")
		if city == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Missing city",
			})
			return
		}
		services.RemoveCityFromUser(userID, city)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
		})

	case "PUT":
		var payload struct {
			Order []string `json:"order"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid JSON",
			})
			return
		}
		services.ReorderCities(userID, payload.Order)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
