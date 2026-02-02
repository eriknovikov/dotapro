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
	if err := ensureDB(DB); err != nil {
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

func ensureDB(DB *bun.DB) error {
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
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if DB.PingContext(ctx) != nil {
		log.Warn().Msg("DB connection lost, attemptying to reconnect...")
		if err := ensureDB(DB); err != nil {
			return fmt.Errorf("failed to reconnect to db: %w", err)
		}
	}
	lastFetchedMatchId, err := fetchLastID(DB)
	if err != nil {
		return fmt.Errorf("err getting last_fetched_match_id: %w", err)
	}
	log.Debug().Int64("LFMD", lastFetchedMatchId).Send()
	matchIds, err := fetchMatchIDs(lastFetchedMatchId)
	if err != nil {
		return fmt.Errorf("err fetching match ids < last_fetched_matched_id: %w", err)
	}
	/*
		10K matches , first from opendota, then into rds.
		array
		repeat 100 times:
			fetch 50 from OD (with retry)
			fetch 50 from OD (with retry)
			insert 100 into rds (with retry)
	*/
	return ScrapeMatches(matchIds, DB)

}
