package handlers

import (
	"net/http"
	service "weatherApp/pkg/services"
)

type Handler struct {
	serv service.Service
}

func NewHandler(serv service.Service) *Handler {
	return &Handler{serv}
}

// api groups
func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/cities", h.getCityList)
	router.HandleFunc("/forecast/short", h.getShortForecast)
	router.HandleFunc("/forecast/details", h.getDetailedForecast)
	return router
}
