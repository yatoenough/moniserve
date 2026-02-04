package app

import (
	"context"
	"log"
	"net/http"

	"github.com/yatoenough/moniserve/internal/handlers"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	mux := http.NewServeMux()

	handlers.Setup(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &App{
		server,
	}
}

func (a *App) Start() error {
	log.Println("Server starting on :8080")
	if err := a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down...")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped gracefully")

	return nil
}
