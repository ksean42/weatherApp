package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"weatherApp/service"
)

type Server struct {
	httpServer *http.Server
	service    *service.Service
	db         *WeatherDB
}

func (s *Server) Start(config *Config, db *WeatherDB) error {
	router := http.NewServeMux()
	s.httpServer = &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.db = db
	router.HandleFunc("/cities", s.getCityList)
	router.HandleFunc("/forecast", s.getForecast)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) getForecast(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request method"))
		return
	}
	args := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(args.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

}
func logs(handler http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		defer handler(w, r)
		body, err := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		data := strings.Builder{}
		data.WriteString(fmt.Sprintf("Request time: %s\nMethod: %s\n", time.Now(), "GET"))
		data.WriteString("\nHeader: \n")
		for k, v := range r.Header {
			data.WriteString(fmt.Sprintf("%s : ", k))
			for _, s := range v {
				data.WriteString(fmt.Sprintf("%s", s))
			}
			data.WriteString("\n")
		}
		data.WriteString("\nQuery params : \n")
		err = r.ParseForm()
		if err == nil {
			for k, v := range r.Form {
				data.WriteString(fmt.Sprintf("%s : ", k))
				for _, s := range v {
					data.WriteString(fmt.Sprintf("%s", s))
				}
				data.WriteString("\n")
			}
		} else {
			data.WriteString(err.Error() + "\n")
		}
		data.WriteString(fmt.Sprintf("Body:\n%s\n", string(body)))
		//l.W.Write([]byte(data.String()))
	}
}

func (s *Server) getCityList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request method"))
		return
	}

	res, err := s.db.GetCityList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
