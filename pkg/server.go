package pkg

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(config *Config, router *http.ServeMux) error {
	s.httpServer = &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// date format: 2022-10-22 11:00:01
//func (s *Server) getDetailedForecast(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Bad request method"))
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	args := r.URL.Query()
//	res := args.Get("time")
//	time, err := time.Parse("2006-01-02 15:04:05", res)
//	id, err := strconv.Atoi(args.Get("id"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Bad request"))
//		return
//	}
//	data, err := s.db.GetDetailedForecast(id, time)
//	resp, err := json.Marshal(data)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Something went wrong"))
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write(resp)
//}
//
//func (s *Server) getForecast(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Bad request method"))
//		return
//	}
//	args := r.URL.Query()
//	w.Header().Set("Content-Type", "application/json")
//	id, err := strconv.Atoi(args.Get("id"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Bad request"))
//		return
//	}
//	data, err := s.db.GetForecast(id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Something went wrong"))
//		return
//	}
//	var resp []byte
//	resp, err = json.Marshal(data)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Something went wrong"))
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write(resp)
//}
//
//func (s *Server) getCityList(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "GET" {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write([]byte("Bad request method"))
//		return
//	}
//
//	res, err := s.db.GetCityList()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Something went wrong"))
//		return
//	}
//	response, err := json.Marshal(res)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte("Something went wrong"))
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(response)
//}
