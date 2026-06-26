package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ServerPort int    `env:"SERVER_PORT" env-default:"8080"`
	ServerHost string `env:"SERVER_HOST" env-default:""`

	LoggerLevel string `env:"LOGGER_LEVEL" env-default:"DEBUG"`

	DbHost     string `env:"DATABASE_HOST"`
	DbPort     int    `env:"DATABASE_PORT"`
	DbUsername string `env:"DATABASE_USERNAME"`
	DbPassword string `env:"DATABASE_PASSWORD"`
	DbName     string `env:"DATABASE_NAME"`
}

func FromEnv() (*Config, error) {
	config := new(Config)
	err := cleanenv.ReadEnv(config)
	return config, err
}
