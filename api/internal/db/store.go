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

func OpenDefault() *Store {
	pgpoolAuth := os.Getenv("PGPOOL_AUTH")
	if pgpoolAuth == "" {
		log.Fatal("PGPOOL_AUTH is missing")
	}

	pgpoolAuthConn, err := pgxpool.New(context.Background(), pgpoolAuth)
	if err != nil {
		log.Fatal("failed to connect to pgpool:", err)
	}

	return &Store{
		db:      pgpoolAuthConn,
		Queries: sqlc.New(pgpoolAuthConn),
	}
}

func OpenTest() *Store {
	pgpoolTestAuth := os.Getenv("PGPOOL_TEST_AUTH")
	if pgpoolTestAuth == "" {
		log.Fatal("PGPOOL_TEST_AUTH is missing")
	}

	pgpoolTestAuthConn, err := pgxpool.New(context.Background(), pgpoolTestAuth)
	if err != nil {
		log.Fatal("failed to connect to pgpool:", err)
	}

	return &Store{
		db:      pgpoolTestAuthConn,
		Queries: sqlc.New(pgpoolTestAuthConn),
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
