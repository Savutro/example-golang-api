package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbUser         string `mapstructure:"DB_USER"`
	DbPassword     string `mapstructure:"DB_PASSWORD"`
	DbHost         string `mapstructure:"DB_HOST"`
	DbPort         string `mapstructure:"DB_PORT"`
	DbName         string `mapstructure:"DB_NAME"`
	CookiePassword string `mapstructure:"COOKIE_PASSWORD"`
	SessionName    string `mapstructure:"SESSION_NAME"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("can't read config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return config, nil
}
