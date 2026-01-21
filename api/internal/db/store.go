package db

import (
	"context"
	"log"
	"os"

	"ykstreaming_api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func Open() *Store {
	postgresAuth := os.Getenv("POSTGRES_AUTH")
	if postgresAuth == "" {
		log.Fatal("POSTGRES_AUTH is missing")
	}

	dbConn, err := pgxpool.New(context.Background(), postgresAuth)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}

	return &Store{
		db:      dbConn,
		Queries: sqlc.New(dbConn),
	}
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Store) Tx(ctx context.Context, handler func(q *sqlc.Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	qtx := s.Queries.WithTx(tx)

	if err := handler(qtx); err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			log.Print(err)
			return rollbackErr
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}
