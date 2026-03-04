package main

import (
	"context"
	"dotapro-lambda-api/config"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/db"
	"fmt"
	"net/http"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/uptrace/bun"
)

// setupDB initializes and returns a database connection pool.
// It retrieves the database URL from either local config or AWS SSM.
func setupDB() (*bun.DB, error) {
	dbURL, err := getDBURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get database URL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBConnectionTimeout)
	defer cancel()

	db, err := db.CreatePool(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	return db, nil
}

// getDBURLFromSSM retrieves the database URL from AWS SSM Parameter Store.
// The parameter is decrypted automatically.
func getDBURLFromSSM(ctx context.Context) (string, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	decryption := true

	output, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &config.CONFIG.DBURLParamName,
		WithDecryption: &decryption,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get SSM parameter: %w", err)
	}

	if output.Parameter.Value == nil {
		return "", fmt.Errorf("SSM parameter value is nil")
	}

	return *output.Parameter.Value, nil
}

// getDBURL returns the database URL based on the current environment.
// For local environment, it uses the LOCAL_DB_URL config.
// For production, it retrieves the URL from AWS SSM Parameter Store.
func getDBURL() (string, error) {
	if config.IsLocal() {
		return config.CONFIG.LocalDBURL, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.SSMTimeout)
	defer cancel()

	url, err := getDBURLFromSSM(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve database URL from SSM: %w", err)
	}

	return url, nil
}

// Middleware to block access from the direct urls (from API GW or CF)
func isRequestAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		secret := r.Header.Get(config.CONFIG.MagicHeaderName)
		if secret != config.CONFIG.MagicHeaderValue {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error": "Direct access forbidden"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
