package db

import (
	"context"
	"dotapro-lambda-api/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// CreatePool creates a new database connection pool with optimized settings for Lambda.
// The pool is configured with conservative limits suitable for Lambda's short-lived execution model.
//
// Pool Configuration (from environment variables or defaults):
// - MaxConns: Maximum number of connections (kept low for Lambda)
// - MinConns: Minimum number of connections to maintain
// - MaxConnIdleTime: Maximum time a connection can be idle before being closed
// - MaxConnLifetime: Maximum lifetime of a connection before being recycled
func CreatePool(ctx context.Context, connStr string) (*bun.DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings from config
	cfg.MaxConns = int32(config.CONFIG.DBMaxConns)
	cfg.MinConns = int32(config.CONFIG.DBMinConns)
	cfg.MaxConnIdleTime = config.CONFIG.DBMaxConnIdleTime
	cfg.MaxConnLifetime = config.CONFIG.DBMaxConnLifetime

	// Create the connection pool
	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Verify the connection is working
	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	// Wrap pgx pool with Bun ORM
	sqldb := stdlib.OpenDBFromPool(p)
	db := bun.NewDB(sqldb, pgdialect.New())
	
	return db, nil
}
