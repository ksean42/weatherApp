package pkg

import (
	"context"
	"net/http"
	"time"
)

//Server is configurable http.server entity
type Server struct {
	httpServer *http.Server
}

//Start server with configuration
func (s *Server) Start(config *Config, router *http.ServeMux) error {
	s.httpServer = &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

//Stop http server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
