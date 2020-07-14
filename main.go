package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Start API Service

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

// Echo is a basic HTTP Handler.
// If you open localhost:8000 in your browser, you may notice
// double requets being made. This happens because the browser
// sends a request in the background for a website favicon.
func Echo(w http.ResponseWriter, r *http.Request) {

	// print writer
	io.WriteString(w, "Hello from a HandleFunc #1!\n")

	// Print a random number at the beginning and end of each request.
	n := rand.Intn(1000)
	start, elapsed := time.Now(), time.Now()
	log.Println(w, "\nstart %s", n)
	defer log.Println("end %s", n)
	log.Println(w, start, "start")
	log.Println(w, elapsed, "elapsed")

	// Simulate a long-running request.
	time.Sleep(3 * time.Second)

	end := elapsed.Sub(start)
	fmt.Fprintln(w, "You asked to %s %s\n", r.Method, r.URL.Path)
	log.Println(w, end, "amount of time elapsed")
}
