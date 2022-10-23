package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"weatherApp/models"
	"weatherApp/pkg"
)

func getCity(city string, apikey string, sw *sync.WaitGroup, dest *models.City) {
	response, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, apikey))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	var data []models.City
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal(err)
	}
	sw.Done()
	*dest = data[0]
}

func initdb() {

}

func main() {
	//start := time.Now()
	config := pkg.NewConfig()
	w := pkg.WeatherDB{}
	w.Connection(config.DBConfig)
	server := new(pkg.Server)

	if err := server.Start(config, &w); err != nil {
		log.Fatal(err)
	}

	var cities []models.City
	//sw := &sync.WaitGroup{}
	var city models.City
	for _, v := range config.Cities {
		//sw.Add(1)
		getCity(v, config.ApiKey, &sync.WaitGroup{}, &city) // multithread
		cities = append(cities, city)
	}
	//sw.Wait()
	//w.SaveCities(cities)
	//w.GetForecast()
}
