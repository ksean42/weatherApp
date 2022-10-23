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

func (w *WeatherDB) GetForecast() {
	res, err := w.DB.Query("select * from cities;")
	if err != nil {
		fmt.Println(err)
		return
	}
	var id, name, con, lon, lat string
	for res.Next() {
		res.Scan(&id, &name, &con, &lon, &lat)
		resp := GetWeather(lon, lat)
		rawData, err := json.Marshal(resp)
		temp := getDayTemp(resp)
		res, err := w.DB.Query("select cityId, date from forecast where cityId=$1;", id)
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.Next() == true {
			_, err = w.DB.Exec("UPDATE forecast SET temp=$1, date=$2, misc=$3 where cityId=$4;",
				temp, resp.List[0].DtTxt, rawData, id)

		} else {
			_, err = w.DB.Exec("INSERT INTO forecast VALUES ($1, $2, $3, $4);",
				id, temp, resp.List[0].DtTxt, rawData)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
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
func (w *WeatherDB) SaveForecast() {

}
