package config

import (
	"strings"

	"github.com/spf13/viper"
)

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

func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	origins := viper.GetString("CORS_ALLOWED_ORIGINS")
	if origins != "" {
		config.AllowedOrigins = strings.Split(origins, ",")
	}

	return config, nil
}
