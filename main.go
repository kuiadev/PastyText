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
	err := startServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v\n", err)
	}
}

func startServer() error {
	pts, err := server.NewPtServer()
	if err != nil {
		return err
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

	return server.Shutdown(ctx)
}
