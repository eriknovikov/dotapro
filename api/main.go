package main

import (
	"context"
	"dotapro-lambda-api/config"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/filtersmetadata"
	"dotapro-lambda-api/matches"
	"dotapro-lambda-api/series"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// Request and Response types for Lambda
type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

// App holds the application state and dependencies
type App struct {
	db                        *bun.DB
	matchModel                *matches.Model
	seriesModel               *series.Model
	filtersMetadataModel      *filtersmetadata.Model
	matchController           *matches.Controller
	seriesController          *series.Controller
	filtersMetadataController *filtersmetadata.Controller
	adapter                   *httpadapter.HandlerAdapterV2
}

// NewApp creates and initializes a new application instance
func NewApp() (*App, error) {
	db, err := setupDB()
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	matchModel := matches.NewModel(db)
	seriesModel := series.NewModel(db)
	filtersMetadataModel := filtersmetadata.NewModel(db)

	return &App{
		db:                        db,
		matchModel:                matchModel,
		seriesModel:               seriesModel,
		filtersMetadataModel:      filtersMetadataModel,
		matchController:           matches.NewController(matchModel),
		seriesController:          series.NewController(seriesModel),
		filtersMetadataController: filtersmetadata.NewController(filtersMetadataModel),
	}, nil
}

// Close cleans up application resources
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// setupRouter configures and returns the HTTP router with all routes and middleware
func (a *App) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	//cors
	allowedOrigins := getAllowedOrigins(config.IsLocal(), config.CONFIG.CLOUDFRONT_URL)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           constants.CORSMaxAge,
	}))
	if config.IsLocal() {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	// routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "hello from home") })
	r.Get("/matches", a.matchController.GetMany)
	r.Get("/matches/{id}", a.matchController.GetOne)
	r.Get("/series", a.seriesController.GetMany)
	r.Get("/series/{id}", a.seriesController.GetOne)
	r.Get("/filtersmetadata/teams", a.filtersMetadataController.SearchTeams)
	r.Get("/filtersmetadata/leagues", a.filtersMetadataController.SearchLeagues)

	return r
}

// ensureDBConnection checks if the database connection is alive and reconnects if needed
func (a *App) ensureDBConnection(ctx context.Context) error {
	if a.matchModel.DB == nil || a.matchModel.DB.PingContext(ctx) != nil {
		db, err := setupDB()
		if err != nil {
			return fmt.Errorf("failed to re-initialize db pool: %w", err)
		}
		a.matchModel.DB.Close()
		a.matchModel.DB = db
		a.seriesModel.DB = db
		a.filtersMetadataModel.DB = db
		a.db = db
	}
	return nil
}

// Handler is the Lambda entry point
func (a *App) Handler(ctx context.Context, req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.LambdaHandlerTimeout)
	defer cancel()

	if err := a.ensureDBConnection(ctx); err != nil {
		return Response{StatusCode: 500}, nil
	}

	return a.adapter.ProxyWithContext(ctx, req)
}

// init initializes the application configuration and logging
func init() {
	if err := config.LoadEnvs(); err != nil {
		panic(fmt.Errorf("failed to load environment variables: %w", err))
	}
	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("invalid configuration: %w", err))
	}
	if config.IsLocal() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	}
}

// main is the application entry point
func main() {
	app, err := NewApp()
	if err != nil {
		panic(fmt.Errorf("failed to initialize application: %w", err))
	}
	defer app.Close()

	r := app.setupRouter()

	if config.IsLocal() {
		if err := http.ListenAndServe(config.CONFIG.LOCAL_ADDR, r); err != nil {
			panic(fmt.Errorf("failed to start server: %w", err))
		}
	} else if config.IsProd() {
		app.adapter = httpadapter.NewV2(r)
		lambda.Start(app.Handler)
	}
}
