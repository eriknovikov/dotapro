package main

import (
	"context"
	"dotapro-lambda-api/config"
	"dotapro-lambda-api/matches"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const DB_CREATION_TIMEOUT = time.Second * 5
const SSM_TIMEOUT = time.Second * 5

type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

var (
	MATCH_MODEL *matches.Model
	ADAPTER     *httpadapter.HandlerAdapterV2
)

func init() {
	if err := config.LoadEnvs(); err != nil {
		log.Fatal().Err(err).Msg("err loading envs")
	}
	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Msg("err invalid config")
	}
	if config.IsLocal() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {
	db, err := setupDB()
	if err != nil {
		log.Fatal().Err(err).Msg("err setting up DB")
	}
	defer db.Close()
	MATCH_MODEL = matches.NewModel(db)
	matchController := matches.NewController(MATCH_MODEL)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "hello from home") })
	r.Get("/matches/{id}", matchController.GetOne)
	if config.IsLocal() {
		log.Info().Str("LOCAL_ADDR", config.CONFIG.LOCAL_ADDR).Msg("running api locally")
		err := http.ListenAndServe(config.CONFIG.LOCAL_ADDR, r)
		if err != nil {
			log.Fatal().Err(err)
		}
	} else if config.IsProd() {
		ADAPTER = httpadapter.NewV2(r)
		lambda.Start(Handler)
	}
}

func Handler(ctx context.Context, req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if MATCH_MODEL.DB == nil || MATCH_MODEL.DB.PingContext(ctx) != nil {
		log.Warn().Msg("DB connection lost, re-initializing db pool...")
		db, err := setupDB()
		if err != nil {
			log.Error().Err(err).Msg("failed to re-initialize db pool")
			return Response{StatusCode: 500}, nil
		}
		MATCH_MODEL.DB.Close()
		MATCH_MODEL.DB = db
	}
	return ADAPTER.ProxyWithContext(ctx, req)

}
