package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// appv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/gen/app/v1"
	appv1connect "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/gen/app/v1/appv1connect"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/handler"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/logging"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/repository"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/service"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/storage"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Setup logging
	logging.Setup(cfg.Environment)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Create error group for managing goroutines
	g, ctx := errgroup.WithContext(ctx)

	// Initialize database
	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer db.Close()

	// Initialize storage
	storage, err := storage.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize storage")
	}

	// Initialize repository
	repo := repository.New(db.Pool)

	// Initialize service
	svc := service.New(repo, storage)

	// Initialize handler
	h := handler.New(svc)

	// Initialize HTTP server with Connect-RPC handler
	mux := http.NewServeMux()
	path, handler := appv1connect.NewApplicationServiceHandler(h)
	mux.Handle(path, handler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Initialize HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Start HTTP server
	g.Go(func() error {
		log.Info().Int("port", cfg.ServerPort).Msg("starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	// Handle shutdown signal
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-sigCh:
			log.Info().Str("signal", sig.String()).Msg("received shutdown signal")
			cancel()
		}
		return nil
	})

	// Handle graceful shutdown
	g.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("shutting down HTTP server")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("server error")
		os.Exit(1)
	}
}
