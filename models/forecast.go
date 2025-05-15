package models

type ForecastItem struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Rain *struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain,omitempty"`
	DtTxt string `json:"dt_txt"`
	Weather []struct {
		Main string `json:"main"`
		Icon string `json:"icon"`
	} `json:"weather"`
}

type ForecastResponse struct {
	List []ForecastItem `json:"list"`
	City struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
	} `json:"city"`
}
