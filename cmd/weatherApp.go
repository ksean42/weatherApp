package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"weatherApp/pkg"
	"weatherApp/pkg/handlers"
	"weatherApp/pkg/repository"
	"weatherApp/pkg/services"
)

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

func main() {
	s := time.Now()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	config := pkg.NewConfig()
	db, err := repository.NewWeatherDB(config.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	serv := service.NewService(db, config, ctx)
	handler := handlers.NewHandler(*serv)
	fmt.Println(time.Until(s))

	server := &pkg.Server{}
	go gracefulShutdown(ctx, cancel, server, exit)
	if err := server.Start(config, handler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc,
	server *pkg.Server, exit chan os.Signal) {
	<-exit
	if err := server.Stop(ctx); err != nil {
		log.Println(err)
	}
	cancel()
	log.Println("Server shouting down...")

}
