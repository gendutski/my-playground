package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a server
	server := &http.Server{Addr: ":8080"}

	// Channel to listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Define a simple handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Run the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Could not listen on :8080: %v\n", err)
		}
	}()
	fmt.Println("Server is ready to handle requests at :8080")

	// Block until a signal is received
	<-stop

	fmt.Println("Shutting down the server...")

	// Create a context with timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Could not gracefully shutdown the server: %v\n", err)
	}

	fmt.Println("Server stopped")
}
