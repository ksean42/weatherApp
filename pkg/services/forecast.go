package service

import (
	"time"
	"weatherApp/pkg/entities"
	"weatherApp/pkg/repository"
)

type ForecastService struct {
	weatherDB repository.Repository
}

func NewForecastService(weatherDB repository.Repository) *ForecastService {
	return &ForecastService{weatherDB}
}

func (f *ForecastService) SaveCities(cities []entities.City) {
	f.weatherDB.SaveCities(cities)
}
func (f *ForecastService) GetShortForecast(id int) (*entities.ShortForecast, error) {
	return f.weatherDB.GetShortForecast(id)
}

func (f *ForecastService) GetDetailedForecast(id int, date time.Time) (*entities.Details, error) {
	return f.weatherDB.GetDetailedForecast(id, date)
}
func (f *ForecastService) GetCityList() ([]entities.CityResponse, error) {
	return f.weatherDB.GetCityList()
}

func (f *ForecastService) SaveForecast(response entities.Forecast, dayTemp float64) error {
	return f.weatherDB.SaveForecast(response, dayTemp)
}
