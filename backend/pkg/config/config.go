package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(ProvideConfig),
)

// environment types for config
type environment string

// Defined environments
const (
	EnvironmentLocal environment = "LOCAL"
	EnvironmentSTAGE environment = "STAGE"
	EnvironmentProd  environment = "PROD"
)

type Config struct {
	Environment environment `envconfig:"ENVIRONMENT" required:"true"`
	HostURL     string      `envconfig:"HOST_URL" required:"true"`
	Database    struct {
		Host     string `envconfig:"DB_HOST" required:"true"`
		Port     string `envconfig:"DB_PORT" required:"true"`
		Username string `envconfig:"DB_USERNAME" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
		Database string `envconfig:"DB_DATABASE" required:"true"`
	}

	GCP struct {
		ProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
	}
}

func ProvideConfig() (Config, error) {

	// load variables from .env file
	if err := godotenv.Load(); err != nil {
		return Config{}, errors.Wrap(err, "failed to load env file")
	}

	// parse env variables to config struct
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse env variables")
	}

	return config, nil
}
