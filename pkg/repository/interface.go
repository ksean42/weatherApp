package repository

import (
	"time"
	"weatherApp/pkg/entities"
)

// Repository to serve storage
type Repository interface {
	SaveCities(cities []entities.City)
	SaveForecast(response entities.Forecast, dayTemp float64) error
	GetShortForecast(id int) (*entities.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*entities.Details, error)
	GetCityList() ([]entities.CityResponse, error)
}
