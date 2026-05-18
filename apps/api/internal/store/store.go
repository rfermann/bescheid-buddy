package store

import (
	"context"
	"fmt"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewStore(dbConfig config.DatabaseConfig) (*Store, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		"disable")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Store{
		Pool:    pool,
		Queries: db.New(pool),
	}, nil
}

func (db *Store) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *Store) Close() {
	db.Pool.Close()
}
