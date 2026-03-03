package main

import (
	"context"
	"fmt"
	"scraper/config"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

const (
	// Database connection pool settings
	DBMaxConns        = 1
	DBMinConns        = 0
	DBMaxConnIdleTime = 30 * time.Minute
	DBMaxConnLifetime = time.Hour
)

// getDBUrlFromSSM retrieves the database URL from AWS SSM Parameter Store
func getDBUrlFromSSM(ctx context.Context) (string, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS default config: %w", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	decryption := true

	output, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &config.CONFIG.DB_URL_PARAM_NAME,
		WithDecryption: &decryption,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get SSM parameter: %w", err)
	}

	if output.Parameter == nil || output.Parameter.Value == nil {
		return "", fmt.Errorf("SSM parameter value is nil")
	}

	return *output.Parameter.Value, nil
}

// getDBUrl retrieves the database URL from the appropriate source
func getDBUrl() (string, error) {
	var dbUrl string
	var err error

	if config.IsLocal() {
		dbUrl = config.CONFIG.LOCAL_DB_URL
		if dbUrl == "" {
			return "", fmt.Errorf("LOCAL_DB_URL is not set")
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), SSM_TIMEOUT)
		defer cancel()

		dbUrl, err = getDBUrlFromSSM(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve DB URL from SSM: %w", err)
		}
	}

	return dbUrl, nil
}

// setupDB creates and initializes a database connection pool
func setupDB() (*bun.DB, error) {
	dbUrl, err := getDBUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_CREATION_TIMEOUT)
	defer cancel()

	db, err := CreatePool(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create DB pool: %w", err)
	}

	return db, nil
}

// CreatePool creates a new database connection pool with pgxpool
func CreatePool(ctx context.Context, connStr string) (*bun.DB, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool
	cfg.MaxConns = DBMaxConns
	cfg.MinConns = DBMinConns
	cfg.MaxConnIdleTime = DBMaxConnIdleTime
	cfg.MaxConnLifetime = DBMaxConnLifetime

	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := p.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqldb := stdlib.OpenDBFromPool(p)
	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}

// fetchLastID retrieves the last fetched match ID from the scraper metadata table
func fetchLastID(db *bun.DB) (int64, error) {
	var res int64
	err := db.QueryRow("SELECT last_fetched_match_id FROM scraper_metadata WHERE id = 1").Scan(&res)
	if err != nil {
		return -1, fmt.Errorf("failed to fetch last match ID: %w", err)
	}
	return res, nil
}

// updateLastID updates the last fetched match ID in the scraper metadata table
func updateLastID(ctx context.Context, db *bun.DB, matchID int64) error {
	result, err := db.NewUpdate().
		Model(&ScraperMetadata{ID: 1}).
		Set("last_fetched_match_id = ?", matchID).
		Where("id = ?", 1).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update last match ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		log.Warn().
			Int64("match_id", matchID).
			Msg("no rows updated when updating last fetched match ID")
	}

	return nil
}
