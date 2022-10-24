package repository

import (
	"time"
	"weatherApp/models"
)

type Repository interface {
	SaveCities(cities []models.City) error
	SaveForecast(response models.Response, id int, dayTemp float64) error
	GetShortForecast(id int) (*models.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*models.Details, error)
	GetCityList() ([]models.City, error)
}
