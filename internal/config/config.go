package config

import (
	"github.com/spf13/viper"
	"strings"
)

// Config represents the application configuration
type Config struct {
	ServerPort         string   `mapstructure:"SERVER_PORT"`
	RedisAddr          string   `mapstructure:"REDIS_ADDR"`
	PostgresDSN        string   `mapstructure:"POSTGRES_DSN"`
	AuthorizerClientID string   `mapstructure:"AUTHORIZER_CLIENT_ID"`
	AuthorizerURL      string   `mapstructure:"AUTHORIZER_URL"`
	RedirectURL        string   `mapstructure:"REDIRECT_URL"`
	ProxyCurlAPIKey    string   `mapstructure:"PROXYCURL_API_KEY"`
	AllowedOrigins     []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

// LoadConfig loads configuration from various sources
func LoadConfig() (Config, error) {
	var config Config

	// Set up Viper
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read configuration from the file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file is not found or not readable
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	// Unmarshal configuration values into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Load allowed origins
	origins := viper.GetString("CORS_ALLOWED_ORIGINS")
	if origins != "" {
		config.AllowedOrigins = strings.Split(origins, ",")
	}

	return config, nil
}
