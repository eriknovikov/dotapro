package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"scraper/config"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

const (
	DB_CREATION_TIMEOUT = 5 * time.Second
	SSM_TIMEOUT         = 5 * time.Second
	LAMBDA_TIMEOUT      = 4 * time.Minute // Leave 1 minute buffer for cleanup
)

var DB *bun.DB

func initialize() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if err := config.LoadEnvs(); err != nil {
		log.Fatal().Err(err).Msg("failed to load environment variables")
	}
	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Msg("configuration validation failed")
	}
	if config.IsLocal() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Kitchen,
		})
	}
}

func main() {
	initialize()

	if err := ensureDB(); err != nil {
		log.Fatal().Err(err).Msg("failed to create or restart database connection")
	}
	defer func() {
		if DB != nil {
			if err := DB.DB.Close(); err != nil {
				log.Error().Err(err).Msg("error closing database connection")
			}
			if err := DB.Close(); err != nil {
				log.Error().Err(err).Msg("error closing bun database")
			}
		}
	}()

	if config.IsProd() {
		lambda.Start(handler)
	} else if config.IsLocal() {
		// Run handler once (aka: invoke the 'lambda' locally)
		ctx, cancel := context.WithTimeout(context.Background(), LAMBDA_TIMEOUT)
		defer cancel()

		err := handler(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("handler execution failed")
		}
	}
}

// ensureDB ensures a valid database connection exists
func ensureDB() error {
	newDB, err := setupDB()
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Error().Err(err).Msg("error closing old database connection")
		}
	}
	DB = newDB
	return nil
}

// handler is the main Lambda handler function
func handler(ctx context.Context) error {
	// Check for context cancellation at the start
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled before handler started: %w", err)
	}

	// Create a context with timeout for DB operations
	dbCtx, cancel := context.WithTimeout(ctx, DB_CREATION_TIMEOUT)
	defer cancel()

	// Check DB connection
	if err := DB.PingContext(dbCtx); err != nil {
		log.Warn().Err(err).Msg("database connection lost, attempting to reconnect...")
		if err := ensureDB(); err != nil {
			return fmt.Errorf("failed to reconnect to database: %w", err)
		}
		log.Info().Msg("successfully reconnected to database")
	}

	// Create a context with timeout for the entire scraping operation
	scrapeCtx, cancel := context.WithTimeout(ctx, LAMBDA_TIMEOUT)
	defer cancel()

	// Execute the scraping
	err := ScrapeMatches(scrapeCtx, DB, config.CONFIG.MAX_BATCHES)
	if err != nil {
		if errors.Is(err, ErrNoNewMatches) {
			log.Info().Msg("no new matches have been inserted since last scrape")
			return nil
		}
		return fmt.Errorf("scraping failed: %w", err)
	}

	log.Info().Msg("scraping completed successfully")
	return nil
}
