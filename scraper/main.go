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

const DB_CREATION_TIMEOUT = time.Second * 5
const SSM_TIMEOUT = time.Second * 5

var DB *bun.DB

func initialize() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if err := config.LoadEnvs(); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Send()
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
		log.Fatal().Err(fmt.Errorf("err creating or restarting db: %w", err)).Send()
	}
	defer DB.DB.Close()
	defer DB.Close()

	if config.IsProd() {
		lambda.Start(handler)
	} else if config.IsLocal() {
		//run handler once (aka: invoke the 'lambda' locally)
		err := handler(context.Background())
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

}

func ensureDB() error {
	newDB, err := setupDB()
	if err != nil {
		return err
	}
	if DB != nil {
		DB.Close()
	}
	DB = newDB
	return nil
}

func handler(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, DB_CREATION_TIMEOUT)
	defer cancel()
	if DB.PingContext(dbCtx) != nil {
		log.Warn().Msg("DB connection lost, attemptying to reconnect...")
		if err := ensureDB(); err != nil {
			return fmt.Errorf("failed to reconnect to db: %w", err)
		}
	}
	err := ScrapeMatches(ctx, DB, config.CONFIG.SCRAPING_LIMIT)
	if err != nil {
		if errors.Is(err, ErrNoNewMatches) {
			log.Warn().Msg("no new matches have been inserted")
			return nil
		}
		return err
	}

	return nil

}
