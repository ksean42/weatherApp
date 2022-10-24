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

func (h *Handler) InitRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/cities", h.getCityList)
	router.HandleFunc("/short", h.getShortForecast)
	router.HandleFunc("/full", h.getDetailedForecast)
	return router
}
