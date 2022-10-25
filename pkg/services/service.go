package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/entities"
	"weatherApp/pkg/repository"
)

type Forecast interface {
	SaveCities(cities []entities.City) error
	SaveForecast(response entities.Forecast, dayTemp float64) error
	GetShortForecast(id int) (*entities.ShortForecast, error)
	GetDetailedForecast(id int, date time.Time) (*entities.Details, error)
	GetCityList() ([]entities.CityResponse, error)
}

type Service struct {
	Forecast
}

func NewService(repo repository.Repository, config *pkg.Config, ctx context.Context) *Service {
	service := &Service{
		NewForecastService(repo),
	}
	service.getApiInfo(config, ctx)

	return service
}

func (s *Service) getApiInfo(config *pkg.Config, ctx context.Context) {
	var cities []entities.City
	var city entities.City
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	st := time.Now()

	for i, v := range config.Cities {
		wg.Add(1)
		go func(id int, v string) {
			getCity(v, config.ApiKey, &city)
			city.Id = id + 1
			mutex.Lock()
			cities = append(cities, city)
			mutex.Unlock()
			wg.Done()
		}(i, v)
	}
	wg.Wait()
	fmt.Println("city: ", time.Until(st))
	for _, v := range cities {
		fmt.Println(v.Id)
	}
	go func() {
		for {
			s.SaveWeather(cities, config)
			fmt.Println("Forecasts updated")
			<-time.After(1 * time.Minute)
		}
	}()
	s.SaveCities(cities)
}

func (s *Service) handleForecasts(forecasts []entities.Forecast) {
	for _, v := range forecasts {
		if s.SaveForecast(v, getDayTemp(v)) != nil {
			continue
		}
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

func (s *Service) SaveWeather(cities []entities.City, config *pkg.Config) {
	var forecasts []entities.Forecast
	var forecast entities.Forecast
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for _, v := range cities {
		wg.Add(1)
		go func(v entities.City) {
			getWeather(&v, &forecast, config.ApiKey)
			forecast.CityId = v.Id
			mutex.Lock()
			forecasts = append(forecasts, forecast)
			mutex.Unlock()
			wg.Done()
		}(v)
	}
	wg.Wait()
	s.handleForecasts(forecasts)
}

func getWeather(city *entities.City, dest *entities.Forecast, apikey string) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=%s",
		city.Lat, city.Lon, apikey)
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
	data.CityId = city.Id
	*dest = data
}
