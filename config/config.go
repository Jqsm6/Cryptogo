package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   `yaml:"server"`
	Postgres `yaml:"postgres"`
}

type Server struct {
	Host         string        `yaml:"Host"`
	Port         string        `yaml:"Port"`
	LoggingLevel string        `yaml:"LoggingLevel"`
	ReadTimeout  time.Duration `yaml:"ReadTimeout"`
	WriteTimeout time.Duration `yaml:"WriteTimeout"`
}

type Postgres struct {
	PostgreSQLHost     string `yaml:"PostgreSQLHost"`
	PostgreSQLPort     string `yaml:"PostgreSQLPort"`
	PostgreSQLUser     string `yaml:"PostgreSQLUser"`
	PostgreSQLPassword string `yaml:"PostgreSQLPassword"`
	PostgreSQLDBName   string `yaml:"PostgreSQLDBName"`
	PgDriver           string `yaml:"PgDriver"`
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config = &Config{}

		err = cleanenv.ReadConfig("config.yml", config)
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}
