package entities

type Details struct {
	CityId  int    `json:"city_id"`
	City    string `json:"city"`
	Details List   `json:"list"`
}

type ShortForecast struct {
	Country   string   `json:"country"`
	City      string   `json:"city"`
	AvgTemp   float64  `json:"avg_temp"`
	DatesList []string `json:"list_of_dates"`
}

type CityResponse struct {
	Id      int     `json:"city_id"`
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lon     float64 `json:"longitude"`
	Lat     float64 `json:"latitude"`
}
