package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/handlers"
	"weatherApp/pkg/middleware"
	"weatherApp/pkg/repository"
	"weatherApp/pkg/services"
)

func main() {
	s := time.Now()
	logger := middleware.NewLogger()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	config := pkg.NewConfig()
	db, err := repository.NewWeatherDB(config.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	serv := service.NewService(ctx, db, config)
	handler := handlers.NewHandler(*serv, logger)
	fmt.Println("Initial  time:", time.Since(s))

	server := &pkg.Server{}
	go gracefulShutdown(ctx, cancel, server, exit)
	if err := server.Start(config, handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc,
	server *pkg.Server, exit chan os.Signal) {
	<-exit
	cancel()
	log.Println("Server shutting down...")
	if err := server.Stop(ctx); err != nil {
		log.Println(err)
	}
}
