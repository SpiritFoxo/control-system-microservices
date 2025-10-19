package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Addr             string
	UsersServiceURL  string
	OrdersServiceURL string
	JWTSecret        string
}

func Load() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	cfg := &Config{
		Addr:             ":" + viper.GetString("GATEWAY_PORT"),
		UsersServiceURL:  viper.GetString("USERS_SERVICE_URL"),
		OrdersServiceURL: viper.GetString("ORDERS_SERVICE_URL"),
		JWTSecret:        viper.GetString("TOKEN_SECRET"),
	}

	if cfg.Addr == ":" || cfg.UsersServiceURL == "" || cfg.OrdersServiceURL == "" {
		log.Fatalf("Missing required configuration: Addr=%s, UsersServiceURL=%s, OrdersServiceURL=%s",
			cfg.Addr, cfg.UsersServiceURL, cfg.OrdersServiceURL)
	}

	return cfg
}
