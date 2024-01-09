package config

import (
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
	EnvironmentProd  environment = "PROD"
)

func (e environment) IsLocal() bool {
	return e == EnvironmentLocal
}

type Config struct {
	Port        string      `envconfig:"PORT" required:"true"`
	Environment environment `envconfig:"ENVIRONMENT" required:"true"`
	Host        string      `envconfig:"HOST" required:"true"`
	AppHost     string      `envconfig:"APP_HOST" required:"true"`
	Database    struct {
		Host     string `envconfig:"DB_HOST" required:"true"`
		Port     string `envconfig:"DB_PORT" required:"true"`
		Username string `envconfig:"DB_USERNAME" required:"true"`
		Password string `envconfig:"DB_PASSWORD" required:"true"`
		Database string `envconfig:"DB_DATABASE" required:"true"`
	}

	GCP struct {
		ProjectID string `envconfig:"GCP_PROJECT_ID" required:"true"`
		Region    string `envconfig:"GCP_REGION"`
	}
}

func ProvideConfig() (Config, error) {

	// parse env variables to config struct
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse env variables")
	}

	return config, nil
}
