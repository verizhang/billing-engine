package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GRPCPort                      string `envconfig:"GRPC_PORT" default:"9090"`
	RESTPort                      string `envconfig:"REST_PORT" default:"80"`
	PostgresHost                  string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresUsername              string `envconfig:"POSTGRES_USERNAME" default:"5432"`
	PostgresPassword              string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresDatabase              string `envconfig:"POSTGRES_DATABASE" default:"admin"`
	PostgresPort                  string `envconfig:"POSTGRES_PORT" default:"postgres"`
	PostgresSslmode               string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
	PostgresTimeZone              string `envconfig:"POSTGRES_TIMEZONE" default:"100"`
	PostgresMaxConnections        int    `envconfig:"POSTGRES_MAX_CONNECTIONS" default:"100"`
	PostgresMaxIdleConnection     int    `envconfig:"POSTGRES_MAX_IDLE_CONNECTIONS" default:"10"`
	PostgresConnectionMaxIdleTime int    `envconfig:"POSTGRES_CONNECTIONS_MAX_IDLE_TIME" default:"3600"`
}

func New() Config {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
