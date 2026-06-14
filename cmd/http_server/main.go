package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/noedaka/go-config-parser/cmd/http_server/internal/handler"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
	"github.com/noedaka/go-config-parser/internal/service/rules"
)

func main() {
	r := chi.NewRouter()

	rules := []service.Rule{
		rules.DebugLogRule{},
		rules.PlaintextPasswordRule{},
		rules.ZeroHostRule{},
		rules.TLSDisabledRule{},
		rules.NewWeakAlgorithmRule(),
	}

	parser := parser.YamlJsonParser{}
	handler := handler.NewHandler(rules, parser)

	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/config", func(r chi.Router) {
				r.Route("/recommendations", func(r chi.Router) {
					r.Post("/file", handler.ConfigRecommendationsByFileHandler)
					r.Post("/", handler.ConfigRecommendationsByBodyHandler)
				})
			})
		})
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		}
	}()

	log.Printf("Server is listening on 8080")
	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	time.Sleep(1 * time.Second)
}