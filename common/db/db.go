package db

import (
	"context"
	"errors"

	"github.com/adryanev/go-http-service-template/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB provides access to the database
type DB struct {
	Pool    *pgxpool.Pool
	Queries *repository.Queries
}

// New creates a new DB instance
func New(pool *pgxpool.Pool, queries *repository.Queries) (*DB, error) {
	if pool == nil {
		return nil, errors.New("cannot use nil database pool")
	}
	if queries == nil {
		return nil, errors.New("cannot use nil queries")
	}
	return &DB{
		Pool:    pool,
		Queries: queries,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping checks if the database connection is alive
func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Legacy global variables (deprecated, will be removed)
var (
	Pool    *pgxpool.Pool
	Queries *repository.Queries
)

// SetDatabase sets the global database pool (deprecated)
func SetDatabase(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("cannot assign nil database")
	}
	Pool = newPool
	return nil
}

// SetQueries sets the global queries (deprecated)
func SetQueries(newQueries *repository.Queries) error {
	if newQueries == nil {
		return errors.New("cannot assign nil queries")
	}
	Queries = newQueries
	return nil
}
