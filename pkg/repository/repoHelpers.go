package repository

import (
	"time"
	"weatherApp/models"
)

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

func getAvgTemp(forecast models.Response) float64 {
	sum := 0.0
	for _, v := range forecast.List {
		sum += v.Main.Temp
	}
	length := float64(len(forecast.List))
	return sum / length
}
