package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/entities"
	"weatherApp/pkg/repository"
)

type Forecast interface {
	SaveCities(cities []entities.City) error
	SaveForecast(response entities.Forecast, id int, dayTemp float64) error
	GetShortForecast(id int) (*entities.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*entities.Details, error)
	GetCityList() ([]entities.City, error)
}

type Service struct {
	Forecast
}

func NewService(repo repository.Repository, config *pkg.Config) *Service {
	service := &Service{
		NewForecastService(repo),
	}
	service.getApiInfo(config)

	return service
}

func (s *Service) getApiInfo(config *pkg.Config) {
	var cities []entities.City
	var forecasts []entities.Forecast
	var forecast entities.Forecast
	var city entities.City
	for _, v := range config.Cities {
		getCity(v, config.ApiKey, &city) // multithread
		cities = append(cities, city)
	}
	for _, v := range cities {
		forecast = getWeather(v.Lon, v.Lat)
		forecasts = append(forecasts, forecast)
	}
	s.SaveCities(cities)

	for i, v := range forecasts {
		s.SaveForecast(v, i+1, getDayTemp(v)) // handle error
	}
}

func getCity(city string, apikey string, dest *entities.City) {
	url := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apikey)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	var data []entities.City
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal(err)
	}
	*dest = data[0]
}

func getDayTemp(resp entities.Forecast) float64 {
	for _, k := range resp.List {
		if h, _, _ := time.Unix(int64(k.Dt), 0).Clock(); h > 12 && h < 16 {
			return k.Main.Temp
		}
	}
	return resp.List[0].Main.Temp
}
func getWeather(lon, lat float64) entities.Forecast {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=0eee4a21ef9a8817b2663009a78009fa", lat, lon)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	var data entities.Forecast
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
