package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adryanev/go-http-service-template/common/db"
	"github.com/adryanev/go-http-service-template/common/messaging"
	"github.com/adryanev/go-http-service-template/repository"

	"github.com/rs/zerolog/log"

	zerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/joho/godotenv"

	_ "github.com/samber/lo"
	_ "github.com/samber/mo"
)

func main() {
	// INITIATE CONFIGURATION
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file, using environment variables")
	}

	cfg := defaultConfig()
	cfg.loadFromEnv()

	// Create a base context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// INITIATE DATABASES
	dbConn, err := setupDatabase(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup database")
	}
	defer dbConn.Close()

	// INITIATE NATS CLIENT
	natsClient, err := setupNatsClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup NATS client")
	}
	defer natsClient.Close()

	// INITIATE SERVER
	server, err := NewAppHttpServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create the server")
	}

	// Inject dependencies
	server.SetDB(dbConn)
	server.SetNatsClient(natsClient)

	// Setup routes
	server.setupRoute()

	// Start server in a goroutine
	go func() {
		if err := server.start(); err != nil {
			log.Error().Err(err).Msg("Server error")
			cancel()
		}
	}()

	log.Info().Msg("Server started successfully")

	// Wait for shutdown signal
	<-shutdown
	log.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.stop(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server gracefully stopped")
}

// setupDatabase initializes the database connection
func setupDatabase(ctx context.Context, cfg config) (*db.DB, error) {
	config, err := pgxpool.ParseConfig(cfg.PgSql.ConnStr())
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	// Setup logger
	logger := zerolog.NewLogger(log.Logger)
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: tracelog.LogLevelInfo,
	}

	pgsqlClient, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	// Test connection
	if err := pgsqlClient.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	queries := repository.New(pgsqlClient)

	// Create DB struct for dependency injection
	dbConn, err := db.New(pgsqlClient, queries)
	if err != nil {
		return nil, fmt.Errorf("creating DB handler: %w", err)
	}

	// Support legacy global access (will be removed in future)
	if err := db.SetDatabase(pgsqlClient); err != nil {
		return nil, fmt.Errorf("setting database: %w", err)
	}

	if err := db.SetQueries(queries); err != nil {
		return nil, fmt.Errorf("setting queries: %w", err)
	}

	return dbConn, nil
}

// setupNatsClient initializes the NATS client
func setupNatsClient(cfg config) (*messaging.NatsClient, error) {
	natsConfig := messaging.Config{
		URL:      cfg.Nats.URL,
		Username: cfg.Nats.Username,
		Password: cfg.Nats.Password,
	}

	client, err := messaging.NewNatsClient(natsConfig)
	if err != nil {
		return nil, fmt.Errorf("creating NATS client: %w", err)
	}

	return client, nil
}
