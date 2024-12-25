package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	ServerPort       string   `mapstructure:"SERVER_PORT"`
	PostgresDSN      string   `mapstructure:"POSTGRES_DSN"`
	AllowedOrigins   []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	EmailFromAddress string   `mapstructure:"EMAIL_FROM_ADDRESS"`
	EmailFromName    string   `mapstructure:"EMAIL_FROM_NAME"`
	SMTPServer       string   `mapstructure:"SMTP_SERVER"`
	SMTPPort         int      `mapstructure:"SMTP_PORT"`
	SMTPUser         string   `mapstructure:"SMTP_USER"`
	SMTPPassword     string   `mapstructure:"SMTP_PASSWORD"`
	SMTPHost         string   `mapstructure:"SMTP_HOST"`
	AppURL           string   `mapstructure:"APP_URL"`
	JWTSecret        string   `mapstructure:"JWT_SECRET"`
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
