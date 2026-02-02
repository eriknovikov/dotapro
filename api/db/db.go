package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func CreatePool(ctx context.Context, connStr string) (*bun.DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 2
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute
	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	sqldb := stdlib.OpenDBFromPool(p)

	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}
