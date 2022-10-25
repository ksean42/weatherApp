package handlers

import (
	"net/http"
	service "weatherApp/pkg/services"
)

//Handler struct
type Handler struct {
	serv service.Service
}

//NewHandler - server constructor
func NewHandler(serv service.Service) *Handler {
	return &Handler{serv}
}

// InitRouter initialize and return router
func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/cities", h.getCityList)
	router.HandleFunc("/forecast/short", h.getShortForecast)
	router.HandleFunc("/forecast/details", h.getDetailedForecast)
	return router
}
