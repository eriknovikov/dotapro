package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvironLocal = "local"
	EnvironProd  = "prod"
)

var CONFIG = &Config{}

type Config struct {
	LocalDBURL      string `env:"LOCAL_DB_URL" env-default:"postgres://postgres:admin@172.17.0.1:15432/dotapro"`
	DBURLParamName  string `env:"DB_URL_PARAM_NAME"`
	Environ         string `env:"ENVIRON" env-default:"prod"`
	MaxBatches      int    `env:"MAX_BATCHES" env-default:"50"`
	// Scraper configuration
	BatchSize       int           `env:"BATCH_SIZE" env-default:"800"`
	MaxRetries      int           `env:"MAX_RETRIES" env-default:"3"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"30s"`
	DBTimeout       time.Duration `env:"DB_TIMEOUT" env-default:"5s"`
}

func LoadEnvs() error {
	return cleanenv.ReadEnv(CONFIG)
}
func Validate() error {
	errs := []string{}
	if CONFIG.Environ != EnvironLocal && CONFIG.Environ != EnvironProd {
		errs = append(errs, "ENVIRON can only be local or prod")
	}
	_, isRunningOnLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	if CONFIG.Environ == EnvironProd && !isRunningOnLambda {
		errs = append(errs, "ENVIRON is prod despite not being within lambda")
		return errors.New(strings.Join(errs, "; "))
	} else if CONFIG.Environ == EnvironLocal && isRunningOnLambda {
		errs = append(errs, "ENVIRON is local despite being within lambda")
		return errors.New(strings.Join(errs, "; "))
	}

	if CONFIG.Environ == EnvironLocal {
		if CONFIG.LocalDBURL == "" {
			errs = append(errs, "LOCAL_DB_URL must be set when ENVIRON=local")
		}

	}
	if CONFIG.Environ == EnvironProd {

		if CONFIG.DBURLParamName == "" {
			errs = append(errs, "DB_URL_PARAM_NAME must be set when ENVIRON=prod")
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "; "))

}

func IsLocal() bool {
	return CONFIG.Environ == EnvironLocal
}
func IsProd() bool {
	return CONFIG.Environ == EnvironProd
}
