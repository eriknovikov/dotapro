package config

import (
	"errors"
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
	LocalDBURL       string `env:"LOCAL_DB_URL" env-default:"postgres://postgres:admin@localhost:5432/dotapro"`
	LocalAddr        string `env:"LOCAL_ADDR" env-default:"http://localhost:8080"`
	DBURLParamName   string `env:"DB_URL_PARAM_NAME"`
	Environ          string `env:"ENVIRON"`
	MagicHeaderName  string `env:"MAGIC_HEADER_NAME"`
	MagicHeaderValue string `env:"MAGIC_HEADER_VALUE"`
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
	if CONFIG.Environ == "" {
		errs = append(errs, "ENVIRON is required (must be 'local' or 'prod')")
	} else if CONFIG.Environ != EnvironLocal && CONFIG.Environ != EnvironProd {
		errs = append(errs, "ENVIRON can only be local or prod")
	}
	if CONFIG.Environ == EnvironLocal {
		if CONFIG.LocalDBURL == "" {
			errs = append(errs, "LOCAL_DB_URL must be set when ENVIRON=local")
		}
		if CONFIG.LocalAddr == "" {
			errs = append(errs, "LOCAL_ADDR must be set when ENVIRON=local")
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
