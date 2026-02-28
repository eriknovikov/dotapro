package config

import (
	"errors"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var CONFIG = &Config{}

type Config struct {
	LOCAL_DB_URL      string `env:"LOCAL_DB_URL" env-default:"postgres://postgres:admin@172.17.0.1:15432/dotapro"`
	LOCAL_ADDR        string `env:"LOCAL_ADDR" env-default:"localhost:8080"`
	DB_URL_PARAM_NAME string `env:"DB_URL_PARAM_NAME"`
	ENVIRON           string `env:"ENVIRON" env-default:"prod"`
	
	// Database pool configuration
	DBMaxConns        int           `env:"DB_MAX_CONNS" env-default:"2"`
	DBMinConns        int           `env:"DB_MIN_CONNS" env-default:"1"`
	DBMaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" env-default:"5m"`
	DBMaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" env-default:"30m"`
}

func LoadEnvs() error {
	return cleanenv.ReadEnv(CONFIG)
}

func Validate() error {
	errs := []string{}
	if CONFIG.ENVIRON != "local" && CONFIG.ENVIRON != "prod" {
		errs = append(errs, "ENVIRON can only be local or prod")
	}
	if CONFIG.ENVIRON == "local" {
		if CONFIG.LOCAL_DB_URL == "" {
			errs = append(errs, "LOCAL_DB_URL must be set when ENVIRON=local")
		}
		if CONFIG.LOCAL_ADDR == "" {
			errs = append(errs, "LOCAL_ADDR must be set when ENVIRON=local")
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
