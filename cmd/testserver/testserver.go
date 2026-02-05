package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	servers := make([]*http.Server, 0, 5)

	for i := range 5 {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			status := http.StatusOK
			if i == 2 {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
		})

		server := http.Server{
			Addr:    fmt.Sprintf(":808%d", i),
			Handler: mux,
		}

		servers = append(servers, &server)
	}

	for _, s := range servers {
		go func() {
			log.Printf("Starting server on %s\n", s.Addr)
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Server error: %v", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for _, s := range servers {
		wg.Go(func() {
			log.Printf("Stopping server on %s\n", s.Addr)
			if err := s.Shutdown(ctx); err != nil {
				log.Printf("Shutdown error: %v", err)
			}
		})
	}

	wg.Wait()

	log.Println("All servers stopped")
}
