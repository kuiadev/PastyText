package main

import (
	"fmt"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	errChan := make(chan error)
	go func(ec chan<- error) {
		ec <- startServer()
	}(errChan)

	select {
	case result := <-errChan:
		if result != nil {
			t.Errorf("Starting server returned error: %v", result)
		}
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout probably means the server is serving successfully")
	}
}
