package main

import (
	"context"
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
	dbCreationTimeout = 5 * time.Second
	ssmTimeout         = 5 * time.Second
	lambdaTimeout      = 4 * time.Minute // Leave 1 minute buffer for cleanup
)

var DB *bun.DB

func initialize() {
	if err := config.LoadEnvs(); err != nil {
		panic(fmt.Errorf("failed to load environment variables: %w", err))
	}
	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("configuration validation failed: %w", err))
	}
	if config.IsLocal() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Kitchen,
		})
	} else {
		// save those cloudwatch pennies !
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func main() {
	initialize()

	if err := ensureDB(); err != nil {
		panic(fmt.Errorf("failed to create or restart database connection: %w", err))
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
		ctx, cancel := context.WithTimeout(context.Background(), lambdaTimeout)
		defer cancel()

		err := handler(ctx)
		if err != nil {
			panic(fmt.Errorf("handler execution failed: %w", err))
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
	dbCtx, cancel := context.WithTimeout(ctx, dbCreationTimeout)
	defer cancel()

	// Check DB connection
	if err := DB.PingContext(dbCtx); err != nil {
		if err := ensureDB(); err != nil {
			return fmt.Errorf("failed to reconnect to database: %w", err)
		}
	}

	// Create a context with timeout for the entire scraping operation
	scrapeCtx, cancel := context.WithTimeout(ctx, lambdaTimeout)
	defer cancel()

	// Execute the scraping
	err := ScrapeMatches(scrapeCtx, DB, config.CONFIG.MaxBatches)
	if err != nil {
		return fmt.Errorf("scraping failed: %w", err)
	}

	return nil
}
