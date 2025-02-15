package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kuiadev/pastytext/server"
)

func main() {
	startServer()
}

func startServer() {
	pts, err := server.NewPtServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v\n", err)
		return
	}

	server := &http.Server{
		Handler:      pts,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	log.Println("Server started")

	errc := make(chan error, 1)
	go func() {
		errc <- server.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("Failed to start server: %v\n", err)
	case sig := <-sigs:
		log.Printf("Signal received: %v\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Shutdown: %v\n", err)
	}
}
