package main

import (
	"account/internal/config"
	"account/internal/logger"
	"account/internal/repository"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed load config: %v", err)
		return
	}

	logger := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
	if err != nil {
		logger.Error().Msgf("Failed to connect to database: %v", err)
		return
	}
	logger.Info().Msg("database connected")

	repo := repository.NewRepository(db, &logger)
	_ = repo

	logger.Info().Msg("service satrting up")
}
