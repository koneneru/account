package main

import (
	_ "account/docs"
	"account/internal/config"
	"account/internal/logger"
	"account/internal/repository"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Account Service
// @version 1.0
// @description Account Service
// @host localhost:9000
// @BasePath /
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

	router := gin.Default()
	err = router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		logger.Error().Msgf("Failed to load trusted proxies: %v", err)
		return
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", PingExample)
	router.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	logger.Info().Msg("service satrting up")
}

// PingExample godoc
// @Summary Check service availability
// @Description Returns pong
// @Tags health
// @Success 200 {string} string "pong"
// @Router /ping [get]
func PingExample(c *gin.Context) {
	c.String(200, "pong")
}
