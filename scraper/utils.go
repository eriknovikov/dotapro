package main

import (
	"context"
	"scraper/config"
	"time"

	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func getDBUrlFromSSM(ctx context.Context) (string, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("err loading aws default config: %w", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	decrpytion := true
	output, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &config.CONFIG.DB_URL_PARAM_NAME,
		WithDecryption: &decrpytion,
	})

	if err != nil {
		return "", fmt.Errorf("err getting ssm parameter: %w", err)
	}
	if output.Parameter.Value == nil {
		return "", fmt.Errorf("parameter value is nil")
	}

	return *output.Parameter.Value, nil
}

func getDBUrl() (string, error) {
	var dbUrl string
	if config.IsLocal() {
		dbUrl = config.CONFIG.LOCAL_DB_URL
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), SSM_TIMEOUT)
		defer cancel()
		url, err := getDBUrlFromSSM(ctx)
		if err != nil {
			return "", fmt.Errorf("err retrieving db url from SSM PS: %w", err)
		}
		dbUrl = url
	}
	return dbUrl, nil
}

func setupDB() (*bun.DB, error) {
	dbUrl, err := getDBUrl()
	if err != nil {
		return nil, fmt.Errorf("err getting db url: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), DB_CREATION_TIMEOUT)
	defer cancel()
	db, err := CreatePool(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating dbpool: %w", err)
	}

	return db, nil
}

func CreatePool(ctx context.Context, connStr string) (*bun.DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 1
	cfg.MinConns = 0
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.MaxConnLifetime = time.Hour
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

func fetchLastID(db *bun.DB) (int64, error) {
	var res int64
	if err := db.QueryRow("SELECT last_fetched_match_id FROM scraper_metadata WHERE id = 1").Scan(&res); err != nil {
		return -1, err
	}
	return res, nil

}
