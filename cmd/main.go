package main

import (
	_ "account/docs"
	"account/internal/config"
	"account/internal/logger"
	"account/internal/repository"
	"account/internal/server"
	"account/internal/service"
	"database/sql"
	"fmt"
	"log"
	"net"

	accountpb "github.com/koneneru/contracts/account/go"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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

	migrationConn, err := db.DB()
	if migrationConn == nil {
		logger.Fatal().Err(err).Msg("failed to select migrations dialect")
	}

	dbGoose, err := sql.Open("postgres", cfg.DbDsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create sql connection")
	}

	if err := goose.Up(dbGoose, "internal/migrations"); err != nil {
		logger.Fatal().Err(err).Msg("failed to run up migrations")
	}

	repo := repository.NewRepository(db, &logger)
	service := service.New(repo, &logger)
	server := server.New(service, &logger)

	s := grpc.NewServer()
	accountpb.RegisterAccountServer(s, server)

	listenAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Error().Msgf("Failed to listen on %s: %v", listenAddr, err)
		return
	}

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthSrv)
	reflection.Register(s)

	logger.Info().Msgf("gRPC server listening on %s", listenAddr)
	if err := s.Serve(lis); err != nil {
		logger.Error().Msgf("Failed to serve gRPC: %v", err)
		return
	}

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
