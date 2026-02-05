package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yatoenough/moniserve/internal/checker"
	"github.com/yatoenough/moniserve/internal/config"
	"github.com/yatoenough/moniserve/internal/handlers"
)

type App struct {
	addr   string
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	mux := http.NewServeMux()

	ch := checker.NewChecker(cfg.Endpoints, time.Second*5)
	status := handlers.NewStatusHandler(ch)

	mux.HandleFunc("/status", status.Handle)

	addr := fmt.Sprintf(":%s", cfg.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &App{
		addr,
		server,
	}
}

func (a *App) Start() error {
	log.Printf("Server starting on %s\n", a.addr)
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
