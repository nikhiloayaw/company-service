package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	CompanyServiceHost string `mapstructure:"COMPANY_SERVICE_HOST"`
	CompanyServicePort string `mapstructure:"COMPANY_SERVICE_PORT"`
}

var envs = []string{
	"COMPANY_SERVICE_HOST", "COMPANY_SERVICE_PORT",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
