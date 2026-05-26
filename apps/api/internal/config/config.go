package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Auth     AuthConfig     `mapstructure:"auth"`
	API      APIConfig      `mapstructure:"api"`
	Database DatabaseConfig `mapstructure:"database"`
	Stripe   StripeConfig   `mapstructure:"stripe"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port             string `mapstructure:"port"`
	Host             string `mapstructure:"host"`
	ReadTimeout      int    `mapstructure:"read_timeout"`
	WriteTimeout     int    `mapstructure:"write_timeout"`
	GracefulShutdown int    `mapstructure:"graceful_shutdown"`
}

type AuthConfig struct {
	AuthServerURL string `mapstructure:"auth_server_url"`
	CookieName    string `mapstructure:"cookie_name"`
}

type APIConfig struct {
	Title   string `mapstructure:"title"`
	Version string `mapstructure:"version"`
}

type StripeConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret"`
	SecretKey     string `mapstructure:"secret_key"`
	Origin        string `mapstructure:"origin"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// LoadConfig loads configuration using Viper
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	// Configure environment variables
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind environment variables to ensure they work even without config files
	v.BindEnv("server.port")
	v.BindEnv("server.host")
	v.BindEnv("server.read_timeout")
	v.BindEnv("server.write_timeout")
	v.BindEnv("server.graceful_shutdown")
	v.BindEnv("api.title")
	v.BindEnv("api.version")
	v.BindEnv("database.host")
	v.BindEnv("database.user")
	v.BindEnv("database.password")
	v.BindEnv("database.port")
	v.BindEnv("database.dbname")
	v.BindEnv("stripe.webhook_secret")
	v.BindEnv("stripe.secret_key")
	v.BindEnv("stripe.origin")

	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateRequired(v); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateRequired checks that required configuration fields are set using Viper's IsSet()
func validateRequired(v *viper.Viper) error {
	var errors []string

	requiredFields := []string{
		"server.port",
		"server.host",
	}

	for _, field := range requiredFields {
		if !v.IsSet(field) {
			errors = append(errors, fmt.Sprintf("%s is required", field))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// setDefaults sets the default configuration values for optional settings only
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.graceful_shutdown", 30)
	v.SetDefault("api.title", "Bored Template")
	v.SetDefault("api.version", "0.0.1")
}

// GetAddr returns the server address in host:port format
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
