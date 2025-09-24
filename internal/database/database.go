package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/ananthvk/gochat/internal/config"
	"github.com/ananthvk/gochat/internal/database/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	Pool    *pgxpool.Pool
	Querier db.Querier
}

func NewDatabaseService(ctx context.Context, cfg *config.Config) (*DatabaseService, error) {
	pool, err := pgxpool.New(ctx, cfg.DbDSN)
	if err != nil {
		slog.Error("error while creating database service", "error", err)
		return nil, err
	}
	querier := db.New(pool)
	dbService := &DatabaseService{Pool: pool, Querier: querier}

	// Ping the database to check if it's online
	err = PingDatabase(dbService, ctx, cfg.DbPingTimeout)
	if err != nil {
		slog.Error("error while creating database service", "error", err)
		return nil, err
	}

	slog.Info("database service created")
	return dbService, nil
}

// PingDatabase establishes a connection with the database. It is used to check that the database is online. If the connection cannot
// be established within timeout duration, this function will return an error
func PingDatabase(db *DatabaseService, ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := db.Pool.Ping(ctx)
	if err != nil {
		slog.Error("database ping failed", "error", err)
		return err
	}
	return nil
}
