package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/entities"
	"weatherApp/pkg/repository"
)

type Forecast interface {
	SaveCities(cities []entities.City)
	SaveForecast(response entities.Forecast, dayTemp float64) error
	GetShortForecast(id int) (*entities.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*entities.Details, error)
	GetCityList() ([]entities.CityResponse, error)
}

type Service struct {
	Forecast
}

func NewService(ctx context.Context, repo repository.Repository, config *pkg.Config) *Service {
	service := &Service{
		NewForecastService(repo),
	}
	service.getAPIInfo(ctx, config)

	return service
}

func (s *Service) getAPIInfo(ctx context.Context, config *pkg.Config) {
	cities := getCityList(config)
	go s.backgroundUpdate(ctx, cities, config)
	s.SaveCities(cities)
}

func getCityList(config *pkg.Config) []entities.City {
	var cities []entities.City
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i, v := range config.Cities {
		wg.Add(1)
		go func(id int, v string) {
			var city entities.City
			getCityFromAPI(v, config.APIKey, &city)
			city.ID = id + 1
			mutex.Lock()
			cities = append(cities, city)
			mutex.Unlock()
			wg.Done()
		}(i, v)
	}
	wg.Wait()

	return cities
}

func (s *Service) backgroundUpdate(ctx context.Context, cities []entities.City, config *pkg.Config) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.saveWeather(cities, config)
			log.Println("Forecasts updated")
		case <-ctx.Done():
			log.Println("Updating stopped..")
			return
		}
	}
}

func getCityFromAPI(city string, apikey string, dest *entities.City) {
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
		if h, _, _ := time.Unix(k.Dt, 0).Clock(); h > 12 && h < 16 {
			return k.Main.Temp
		}
	}
	return resp.List[0].Main.Temp
}

func (s *Service) saveWeather(cities []entities.City, config *pkg.Config) {
	wg := &sync.WaitGroup{}
	for _, v := range cities {
		wg.Add(1)
		go func(v entities.City) {
			var forecast entities.Forecast
			getForecastFromAPI(&v, &forecast, config.APIKey)
			forecast.CityID = v.ID
			s.SaveForecast(forecast, getDayTemp(forecast))
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func getForecastFromAPI(city *entities.City, dest *entities.Forecast, apikey string) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=%s",
		city.Lat, city.Lon, apikey)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	var data entities.Forecast
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal(err)
	}
	data.CityID = city.ID
	*dest = data
}
