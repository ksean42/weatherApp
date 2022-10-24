package service

import (
	"time"
	"weatherApp/models"
	"weatherApp/pkg/repository"
)

type ForecastService struct {
	weatherDB repository.Repository
}

func NewForecastService(weatherDB repository.Repository) *ForecastService {
	return &ForecastService{weatherDB}
}

func (f *ForecastService) SaveCities(cities []models.City) error {
	return f.weatherDB.SaveCities(cities)
}
func (f *ForecastService) GetShortForecast(id int) (*models.ShortForecast, error) {
	return f.weatherDB.GetShortForecast(id)
}

func (f *ForecastService) GetDetailedForecast(id int, date time.Time) (*models.Details, error) {
	return f.weatherDB.GetDetailedForecast(id, date)
}
func (f *ForecastService) GetCityList() ([]models.City, error) {
	return f.weatherDB.GetCityList()
}

func (f *ForecastService) SaveForecast(response models.Response, id int, dayTemp float64) error {
	return f.weatherDB.SaveForecast(response, id, dayTemp)
}
