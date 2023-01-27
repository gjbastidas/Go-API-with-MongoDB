package env

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DbUsername string `envconfig:"DB_USERNAME" required:"true"`
	DbPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DbHost     string `envconfig:"DB_HOST" required:"true"`
	DbPort     string `envconfig:"DB_PORT" required:"true"`
}

func Config() (*AppConfig, error) {
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		_ = envconfig.Usage("", &conf)
		return nil, err
	}
	return &conf, nil
}
