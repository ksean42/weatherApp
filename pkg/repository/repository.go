package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
	"weatherApp/models"
	"weatherApp/pkg"
)

type WeatherDB struct {
	DB *sqlx.DB
}

func NewDatabase(config *pkg.DBConfig) (*WeatherDB, error) {
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

func (w *WeatherDB) SaveCities(cities []models.City) error {
	for _, v := range cities {
		_, err := w.DB.Exec("INSERT INTO cities(name, country, longitude, latitude) VALUES ($1, $2, $3,$4);",
			v.Name, v.Country, v.Lon, v.Lat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WeatherDB) SaveForecast(forecast models.Response, id int, dayTemp float64) error {
	data, err := json.Marshal(forecast)
	if err != nil {
		return err
	}
	_, err = w.DB.Exec("INSERT INTO forecast VALUES ($1, $2, $3, $4)",
		id, dayTemp, forecast.List[0].DtTxt, data)
	if err != nil {
		return err
	}
	return nil
}

func (w *WeatherDB) GetCityList() ([]models.City, error) {
	var result []models.City

	res, err := w.DB.Query("select * from cities order by name asc;")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var city models.City
		err := res.Scan(&city.Id, &city.Name, &city.Country, &city.Lon, &city.Lat)
		if err != nil {
			return nil, err
		}
		result = append(result, city)
	}
	return result, nil
}

func (w *WeatherDB) GetShortForecast(id int) (*models.ShortForecast, error) {
	var data []byte
	err := w.DB.QueryRow("select misc from forecast where cityId=$1;", id).Scan(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var forecast models.Response
	if err = json.Unmarshal(data, &forecast); err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := &models.ShortForecast{
		Country:   forecast.City.Country,
		City:      forecast.City.Name,
		AvgTemp:   getAvgTemp(forecast),
		DatesList: getDatesList(forecast),
	}
	return response, nil
}

func (w *WeatherDB) GetDetailedForecast(id int, date time.Time) (*models.Details, error) {
	var data []byte
	err := w.DB.QueryRow("select misc from forecast where cityId=$1;", id).Scan(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var forecast models.Response
	if err = json.Unmarshal(data, &forecast); err != nil {
		fmt.Println(err)
		return nil, err
	}
	var resp models.Details
	for i := 1; i < len(forecast.List); i++ {
		curr := forecast.List[i-1].Dt
		next := forecast.List[i].Dt
		if date.Unix() >= curr && date.Unix() < next {
			resp = forecast.List[i-1]
			break
		} else if date.Unix() == next {
			resp = forecast.List[i]
			break
		}
	}
	return &resp, nil
}
