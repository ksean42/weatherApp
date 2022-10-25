package entities

// Details is detailed forecast response entity
type Details struct {
	CityID  int    `json:"city_id"`
	City    string `json:"city"`
	Details List   `json:"list"`
}

// ShortForecast is short forecast response entity
type ShortForecast struct {
	Country   string   `json:"country"`
	City      string   `json:"city"`
	AvgTemp   float64  `json:"avg_temp"`
	DatesList []string `json:"list_of_dates"`
}

// CityResponse is entity for represent city for response
type CityResponse struct {
	ID      int     `json:"city_id"`
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lon     float64 `json:"longitude"`
	Lat     float64 `json:"latitude"`
}
