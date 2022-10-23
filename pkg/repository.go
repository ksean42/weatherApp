package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"weatherApp/models"
)

type Repository interface {
	SaveCities([]string)
	GetForecast()
	SaveForecast()
}

type WeatherDB struct {
	DB *sqlx.DB
}

func (w *WeatherDB) Connection(config *DBConfig) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	w.DB = db
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

func (w *WeatherDB) GetForecast(id int) (*models.ShortForecast, error) {
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
func getAvgTemp(forecast models.Response) float64 {
	sum := 0.0
	for _, v := range forecast.List {
		sum += v.Main.Temp
	}
	length := float64(len(forecast.List))
	return sum / length
}

func getDatesList(forecast models.Response) []string {
	list := make([]string, 0, 6)
	var prev string
	for _, v := range forecast.List {
		t := time.Unix(int64(v.Dt), 0)
		date := t.Format("02-01-2006")

		if date != prev {
			list = append(list, date)
		}
		prev = date
	}
	return list
}

func getDayTemp(resp models.Response) float64 {
	for _, k := range resp.List {
		if h, _, _ := time.Unix(int64(k.Dt), 0).Clock(); h > 12 && h < 16 {
			return k.Main.Temp
		}
	}
	return resp.List[0].Main.Temp
}
func GetWeather(lon, lat string) models.Response {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%s&lon=%s&appid=0eee4a21ef9a8817b2663009a78009fa", lat, lon)
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	var data models.Response
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
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
