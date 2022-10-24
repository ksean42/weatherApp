package repository

import (
	"time"
	"weatherApp/pkg/entities"
)

type Repository interface {
	SaveCities(cities []entities.City) error
	SaveForecast(response entities.Forecast, id int, dayTemp float64) error
	GetShortForecast(id int) (*entities.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*entities.Details, error)
	GetCityList() ([]entities.City, error)
}
