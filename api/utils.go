package main

import (
	"context"
	"dotapro-lambda-api/config"
	"dotapro-lambda-api/db"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/uptrace/bun"
)

func setupDB() (*bun.DB, error) {
	dbUrl, err := getDBUrl()
	if err != nil {
		return nil, fmt.Errorf("err getting db url: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), DB_CREATION_TIMEOUT)
	defer cancel()
	db, err := db.CreatePool(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating dbpool: %w", err)
	}

	return db, nil
}

func getDBUrlFromSSM(ctx context.Context) (string, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("err loading aws default config: %w", err)
	}

	ssmClient := ssm.NewFromConfig(cfg)
	decrpytion := true
	output, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &config.CONFIG.DB_URL_PARAM_NAME,
		WithDecryption: &decrpytion,
	})

	if err != nil {
		return "", fmt.Errorf("err getting ssm parameter: %w", err)
	}
	if output.Parameter.Value == nil {
		return "", fmt.Errorf("parameter value is nil")
	}

	return *output.Parameter.Value, nil
}

func getDBUrl() (string, error) {
	var dbUrl string
	if config.IsLocal() {
		dbUrl = config.CONFIG.LOCAL_DB_URL
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), SSM_TIMEOUT)
		defer cancel()
		url, err := getDBUrlFromSSM(ctx)
		if err != nil {
			return "", fmt.Errorf("err retrieving db url from SSM PS: %w", err)
		}
		dbUrl = url
	}
	return dbUrl, nil
}
