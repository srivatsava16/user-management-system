package database

import (
	"user-management-system/internal/models"

	"github.com/rs/zerolog/log"
)

func RunMigrations() error {
	log.Info().Msg("Running database migrations...")

	err := DB.AutoMigrate(
		&models.BusinessUnit{},
		&models.Division{},
		&models.Role{},
		&models.ModulePermissions{},
		&models.User{},
		&models.AuditLog{},
		&models.PasswordResetToken{},
		&models.TokenInvalidation{},
	)

	if err != nil {
		log.Error().Err(err).Msg("Database migration failed")
		return err
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}
