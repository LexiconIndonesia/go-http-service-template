package main

import (
	"context"
	"lexicon/go-template/common/db"
	"lexicon/go-template/repository"

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
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}
	cfg := defaultConfig()
	cfg.loadFromEnv()

	ctx := context.Background()

	// INITIATE DATABASES
	// PGSQL
	config, err := pgxpool.ParseConfig(cfg.PgSql.ConnStr())

	// logger
	logger := zerolog.NewLogger(log.Logger)

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   logger,
		LogLevel: tracelog.LogLevelInfo,
	}
	pgsqlClient, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to PGSQL Database")
	}
	defer pgsqlClient.Close()

	db.SetDatabase(pgsqlClient)
	queries := repository.New(pgsqlClient)
	db.SetQueries(queries)

	// INITIATE SERVER
	server, err := NewAppHttpServer(cfg)

	if err != nil {
		log.Error().Err(err).Msg("Failed to start the server")
	}

	server.setupRoute()
	server.start()
}
