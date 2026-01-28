package main

import (
	"os"
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
	err = router.Run(":" + port)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
