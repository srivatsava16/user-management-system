package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-system/internal/routes"
	"user-management-system/internal/shared"
	"user-management-system/pkg/database"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	est, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Error().Err(err).Msg("Error loading timezone")
	}
	time.Local = est
}

func main() {

	shared.InitLogger()

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msgf("Could not load .env file, using environment variables: %v", err)
	}

	err = database.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
		return
	}

	err = database.RunMigrations()
	if err != nil {
		log.Fatal().Err(err).Msg("Database migration failed")
	}

	defer database.CloseDB()

	router := routes.SetupRoutes()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8083"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}
