package config

import (
	"HospitalRecord/app/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	HTTP struct {
		IP           string `yaml:"ip" env:"HTTP-IP"`
		Host         string `yaml:"host" env:"HTTP-HOST"`
		Port         string `yaml:"port" env:"HTTP-PORT"`
		ReadTimeout  int    `yaml:"read_timeout" env:"HTTP-READ-TIMEOUT"`
		WriteTimeout int    `yaml:"write_timeout" env:"HTTP-WRITE-TIMEOUT"`
	} `yaml:"http"`
	PostgreSQL struct {
		Username          string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
		Password          string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
		Host              string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
		Port              string `yaml:"port" env:"PSQL_PORT" env-required:"true"`
		Database          string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
		RequestTimeout    int    `yaml:"request_timeout" env-default:"5"`
		ConnectionTimeout int    `yaml:"connection_timeout" env-default:"10"`
		ShutdownTimeout   int    `yaml:"shutdown_timeout" env-default:"5"`
	} `yaml:"postgresql"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	logger := logger.GetLogger()

	once.Do(func() {
		logger.Info("reading the application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	logger.Info("done reading application config")

	return instance
}
