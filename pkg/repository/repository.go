package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/entities"
)

type WeatherDB struct {
	DB *sqlx.DB
}

func NewWeatherDB(config *pkg.DBConfig) (*WeatherDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &WeatherDB{db}, nil
}

func (w *WeatherDB) SaveCities(cities []entities.City) error {
	for _, v := range cities {
		_, err := w.DB.Exec("INSERT INTO cities(city_id, name, country, longitude, latitude) VALUES ($1, $2, $3,$4, $5);",
			v.ID, v.Name, v.Country, v.Lon, v.Lat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WeatherDB) SaveForecast(forecast entities.Forecast, dayTemp float64) error {
	data, err := json.Marshal(forecast)
	if err != nil {
		return err
	}
	var exist bool
	err = w.DB.QueryRow("SELECT exists(SELECT * from forecast where city_id=$1)", forecast.CityID).Scan(&exist)
	if err != nil {
		return err
	}
	if exist {
		_, err = w.DB.Exec("UPDATE forecast SET temp=$1, date=$2, misc = $3  WHERE city_id = $4;",
			dayTemp, forecast.List[0].DtTxt, data, forecast.CityID)
	} else {
		_, err = w.DB.Exec("INSERT INTO forecast VALUES ($1, $2, $3, $4)",
			forecast.CityID, dayTemp, forecast.List[0].DtTxt, data)
	}
	return err
}

func (w *WeatherDB) GetCityList() ([]entities.CityResponse, error) {
	var result []entities.CityResponse

	res, err := w.DB.Query("select * from cities order by name asc;")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var city entities.CityResponse
		err := res.Scan(&city.ID, &city.Name, &city.Country, &city.Lon, &city.Lat)
		if err != nil {
			return nil, err
		}
		result = append(result, city)
	}
	return result, nil
}

func (w *WeatherDB) GetShortForecast(id int) (*entities.ShortForecast, error) {
	var data []byte
	err := w.DB.QueryRow("select misc from forecast where city_id=$1;", id).Scan(&data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var cityName string
	err = w.DB.QueryRow("select name from cities where city_id=$1;", id).Scan(&cityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var forecast entities.Forecast
	if err = json.Unmarshal(data, &forecast); err != nil {
		log.Println(err)
		return nil, err
	}
	response := &entities.ShortForecast{
		Country:   forecast.City.Country,
		City:      cityName,
		AvgTemp:   getAvgTemp(forecast),
		DatesList: getDatesList(forecast),
	}
	return response, nil
}

func (w *WeatherDB) GetDetailedForecast(id int, date time.Time) (*entities.Details, error) {
	var data []byte
	err := w.DB.QueryRow("select misc from forecast where city_id=$1;", id).Scan(&data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var cityName string
	err = w.DB.QueryRow("select name from cities where city_id=$1;", id).Scan(&cityName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var forecast entities.Forecast
	if err = json.Unmarshal(data, &forecast); err != nil {
		log.Println(err)
		return nil, err
	}
	var resp entities.Details
	for i := 1; i < len(forecast.List); i++ {
		curr := forecast.List[i-1].Dt
		next := forecast.List[i].Dt
		if date.Unix() >= curr && date.Unix() < next {
			resp.Details = forecast.List[i-1]
			break
		} else if date.Unix() == next {
			resp.Details = forecast.List[i]
			break
		}
	}
	resp.City = cityName
	resp.CityID = id
	return &resp, nil
}
