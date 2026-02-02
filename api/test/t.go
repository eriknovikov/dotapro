package main

import (
	"context"
	"dotapro-lambda-api/db"
	"dotapro-lambda-api/matches"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
	})
	log.Info().Msg("hi")
	dsn := "postgres://postgres:admin@172.17.0.1:15432/dotapro"
	db, err := db.CreatePool(context.Background(), dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("err setting up DB")
	}

	matchModel := matches.NewModel(db)
	m, err := matchModel.GetOne(context.Background(), 32)
	log.Err(err).Send()
	log.Info().Any("match", m).Send()
}
