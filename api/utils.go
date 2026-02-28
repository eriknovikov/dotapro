package main

import (
	"context"
	"dotapro-lambda-api/config"
	"dotapro-lambda-api/constants"
	"dotapro-lambda-api/db"
	"fmt"

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
		Name:           &config.CONFIG.DB_URL_PARAM_NAME,
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
		return config.CONFIG.LOCAL_DB_URL, nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), constants.SSMTimeout)
	defer cancel()
	
	url, err := getDBURLFromSSM(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve database URL from SSM: %w", err)
	}
	
	return url, nil
}
