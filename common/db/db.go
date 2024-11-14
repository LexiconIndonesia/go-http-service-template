package db

import (
	"errors"
	"lexicon/go-template/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool    *pgxpool.Pool
	Queries *repository.Queries
)

func SetDatabase(newPool *pgxpool.Pool) error {

	if newPool == nil {
		return errors.New("cannot assign nil database")
	}
	Pool = newPool
	return nil
}
func SetQueries(newQueries *repository.Queries) error {
	if newQueries == nil {
		return errors.New("cannot assign nil queries")
	}
	Queries = newQueries
	return nil
}
