package db

import (
	"context"
	"log"
	"os"

	"ykstreaming_api/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	dbPrimary *pgxpool.Pool
	dbReplica *pgxpool.Pool
	WQueries  *sqlc.Queries
	RQueries  *sqlc.Queries
}

func OpenDefault() *Store {
	postgresPrimaryAuth := os.Getenv("POSTGRES_PRIMARY_AUTH")
	if postgresPrimaryAuth == "" {
		log.Fatal("POSTGRES_PRIMARY_AUTH is missing")
	}

	postgresReplicaAuth := os.Getenv("POSTGRES_READ_REPLICA_AUTH")
	if postgresReplicaAuth == "" {
		log.Fatal("POSTGRES_PRIMARY_AUTH is missing")
	}

	dbPrimaryConn, err := pgxpool.New(context.Background(), postgresPrimaryAuth)
	if err != nil {
		log.Fatal("failed to connect to primary DB:", err)
	}

	dbReplicaConn, err := pgxpool.New(context.Background(), postgresPrimaryAuth)
	if err != nil {
		log.Fatal("failed to connect to replica DB:", err)
	}

	return &Store{
		dbPrimary: dbPrimaryConn,
		dbReplica: dbReplicaConn,
		WQueries:  sqlc.New(dbPrimaryConn),
		RQueries:  sqlc.New(dbReplicaConn),
	}
}

func OpenTest() *Store {
	postgresPrimaryAuth := os.Getenv("POSTGRES_PRIMARY_TEST_AUTH")
	if postgresPrimaryAuth == "" {
		log.Fatal("POSTGRES_PRIMARY_TEST_AUTH is missing")
	}

	postgresReplicaAuth := os.Getenv("POSTGRES_READ_REPLICA_TEST_AUTH")
	if postgresReplicaAuth == "" {
		log.Fatal("POSTGRES_READ_REPLICA_TEST_AUTH is missing")
	}

	dbPrimaryConn, err := pgxpool.New(context.Background(), postgresPrimaryAuth)
	if err != nil {
		log.Fatal("failed to connect to primary test DB:", err)
	}

	dbReplicaConn, err := pgxpool.New(context.Background(), postgresPrimaryAuth)
	if err != nil {
		log.Fatal("failed to connect to replica test DB:", err)
	}

	return &Store{
		dbPrimary: dbPrimaryConn,
		dbReplica: dbReplicaConn,
		WQueries:  sqlc.New(dbPrimaryConn),
		RQueries:  sqlc.New(dbReplicaConn),
	}
}

func (s *Store) Close() {
	if s.dbPrimary != nil {
		s.dbPrimary.Close()
	}

	if s.dbReplica != nil {
		s.dbReplica.Close()
	}
}

func (s *Store) Tx(ctx context.Context, handler func(q *sqlc.Queries) error) error {
	tx, err := s.dbPrimary.Begin(ctx)
	if err != nil {
		return err
	}

	qtx := s.WQueries.WithTx(tx)

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
