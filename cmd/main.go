package main

import (
	_ "account/docs"
	"account/internal/app"
	"account/internal/config"
	"account/internal/logger"
	"context"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:9000
// @BasePath /
func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	l := logger.New(cfg)

	application := app.New(&l, cfg)

	if err := application.Run(ctx); err != nil {
		l.Fatal().Err(err).Msg("error")
	}
}
