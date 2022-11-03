package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"weatherApp/pkg/entities"
)

// date format:
// /forecast/details?id=2&time=2022-10-22 18:00
func (h *Handler) getDetailedForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		writeError(w, http.StatusBadRequest, "Bad request method",
			errors.New("bad request method"))
		return
	}
	args := r.URL.Query()
	res := args.Get("time")
	t, err := time.Parse("2006-01-02 15:04", res)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	id, err := strconv.Atoi(args.Get("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	data, err := h.serv.Forecast.GetDetailedForecast(id, t)
	if err != nil {
		writeError(w, http.StatusNotFound, "Data not found", err)
		return
	}
	resp, err := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) getShortForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		writeError(w, http.StatusBadRequest, "Bad request method",
			errors.New("bad request method"))
		return
	}
	args := r.URL.Query()

	id, err := strconv.Atoi(args.Get("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Bad request", err)
		return
	}
	data, err := h.serv.Forecast.GetShortForecast(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Data not found", err)
		return
	}
	var resp []byte
	resp, err = json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) getCityList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		writeError(w, http.StatusBadRequest, "Bad request method",
			errors.New("bad request method"))
		return
	}
	res, err := h.serv.Forecast.GetCityList()
	if err != nil {
		writeError(w, http.StatusNotFound, "Data not found", err)
		return
	}
	response, err := json.Marshal(res)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func writeError(w http.ResponseWriter, status int, message string, err error) {
	log.Println(err)
	w.WriteHeader(status)
	e := &entities.Error{Error: message}
	byteMessage, _ := json.Marshal(e)
	w.Write(byteMessage)
}
