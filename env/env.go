package env

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DbUsername string `envconfig:"DB_USERNAME" default:"admin" required:"true"`
	DbPassword string `envconfig:"DB_PASSWORD" default:"secret" required:"true"`
	DbHost     string `envconfig:"DB_HOST" default:"localhost" required:"true"`
	DbPort     string `envconfig:"DB_PORT" default:"27017" required:"true"`
	DbName     string `envconfig:"DB_NAME" default:"simple-api-with-mongodb" required:"true"`
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
