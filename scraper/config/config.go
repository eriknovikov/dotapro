package config

import (
	"errors"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

var CONFIG = &Config{}

type Config struct {
	LOCAL_DB_URL      string `env:"LOCAL_DB_URL" env-default:"postgres://postgres:admin@172.17.0.1:15432/dotapro"`
	DB_URL_PARAM_NAME string `env:"DB_URL_PARAM_NAME"`
	ENVIRON           string `env:"ENVIRON" env-default:"prod"`
	SCRAPING_LIMIT    int    `env:"SCRAPING_LIMIT" env-default:"800"`
}

func LoadEnvs() error {
	return cleanenv.ReadEnv(CONFIG)
}
func Validate() error {
	errs := []string{}
	if CONFIG.ENVIRON != "local" && CONFIG.ENVIRON != "prod" {
		errs = append(errs, "ENVIRON can only be local or prod")
	}
	_, isRunningOnLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	if CONFIG.ENVIRON == "prod" && !isRunningOnLambda {
		errs = append(errs, "ENVIRON is prod despite not being within lambda")
		return errors.New(strings.Join(errs, "; "))
	} else if CONFIG.ENVIRON == "local" && isRunningOnLambda {
		errs = append(errs, "ENVIRON is local despite being within lambda")
		return errors.New(strings.Join(errs, "; "))
	}

	if CONFIG.ENVIRON == "local" {
		if CONFIG.LOCAL_DB_URL == "" {
			errs = append(errs, "LOCAL_DB_URL must be set when ENVIRON=local")
		}

	}
	if CONFIG.ENVIRON == "prod" {

		if CONFIG.DB_URL_PARAM_NAME == "" {
			errs = append(errs, "DB_URL_PARAM_NAME must be set when ENVIRON=prod")
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "; "))

}

func IsLocal() bool {
	return CONFIG.ENVIRON == "local"
}
func IsProd() bool {
	return CONFIG.ENVIRON == "prod"
}
