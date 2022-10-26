package handlers

import (
	"net/http"
	"weatherApp/pkg/middleware"
	service "weatherApp/pkg/services"
)

// Handler struct
type Handler struct {
	serv service.Service
	*middleware.Logger
}

// NewHandler - server constructor
func NewHandler(serv service.Service, log *middleware.Logger) *Handler {
	return &Handler{serv, log}
}

// InitRouter initialize and return router
func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/cities", h.Log(h.getCityList))
	router.HandleFunc("/forecast/short", h.Log(h.getShortForecast))
	router.HandleFunc("/forecast/details", h.Log(h.getDetailedForecast))
	return router
}
