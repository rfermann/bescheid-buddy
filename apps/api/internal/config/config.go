package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Auth        AuthConfig        `mapstructure:"auth"`
	API         APIConfig         `mapstructure:"api"`
	Database    DatabaseConfig    `mapstructure:"database"`
	S3          S3Config          `mapstructure:"s3"`
	Stripe      StripeConfig      `mapstructure:"stripe"`
	SMTP        SMTPConfig        `mapstructure:"smtp"`
	Vision      VisionConfig      `mapstructure:"vision"`
	Upload      UploadConfig      `mapstructure:"upload"`
	PaperLetter PaperLetterConfig `mapstructure:"paper_letter"`
	Pingen      PingenConfig      `mapstructure:"pingen"`
	LLM         LLMConfig         `mapstructure:"llm"`
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

type S3Config struct {
	Endpoint    string `mapstructure:"endpoint"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	Region      string `mapstructure:"region"`
	FrontendURL string `mapstructure:"frontend_url"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

type VisionConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type UploadConfig struct {
	MaxBytes int64 `mapstructure:"max_bytes"`
}

type PaperLetterConfig struct {
	MaxPDFBytes int64 `mapstructure:"max_pdf_bytes"`
}

type PingenConfig struct {
	ClientID        string `mapstructure:"client_id"`
	ClientSecret    string `mapstructure:"client_secret"`
	OrganisationID  string `mapstructure:"organisation_id"`
	IdentityURL     string `mapstructure:"identity_url"`
	APIURL          string `mapstructure:"api_url"`
	Scope           string `mapstructure:"scope"`
	Product         string `mapstructure:"product"`
	AddressPosition string `mapstructure:"address_position"`
	PrintMode       string `mapstructure:"print_mode"`
	PrintSpectrum   string `mapstructure:"print_spectrum"`
}

type LLMConfig struct {
	OpenAIAPIKey    string `mapstructure:"openai_api_key"`
	GeminiAPIKey    string `mapstructure:"gemini_api_key"`
	AnthropicAPIKey string `mapstructure:"anthropic_api_key"`
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
	v.BindEnv("s3.endpoint")
	v.BindEnv("s3.access_key")
	v.BindEnv("s3.secret_key")
	v.BindEnv("s3.region")
	v.BindEnv("s3.frontend_url")
	v.BindEnv("stripe.webhook_secret")
	v.BindEnv("stripe.secret_key")
	v.BindEnv("stripe.origin")
	v.BindEnv("smtp.host")
	v.BindEnv("smtp.port")
	v.BindEnv("smtp.username")
	v.BindEnv("smtp.password")
	v.BindEnv("smtp.use_ssl")
	v.BindEnv("vision.api_key")
	v.BindEnv("upload.max_bytes")
	v.BindEnv("paper_letter.max_pdf_bytes")
	v.BindEnv("pingen.client_id")
	v.BindEnv("pingen.client_secret")
	v.BindEnv("pingen.organisation_id")
	v.BindEnv("pingen.identity_url")
	v.BindEnv("pingen.api_url")
	v.BindEnv("pingen.scope")
	v.BindEnv("pingen.product")
	v.BindEnv("pingen.address_position")
	v.BindEnv("pingen.print_mode")
	v.BindEnv("pingen.print_spectrum")
	v.BindEnv("llm.openai_api_key")
	v.BindEnv("llm.gemini_api_key")
	v.BindEnv("llm.anthropic_api_key")

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}
	} else if environment := os.Getenv("APP_ENV"); environment != "" {
		v.SetConfigName("config." + environment)
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
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
	v.SetDefault("upload.max_bytes", 15728640)
	v.SetDefault("paper_letter.max_pdf_bytes", 15728640)
	v.SetDefault("pingen.identity_url", "https://identity-staging.pingen.com")
	v.SetDefault("pingen.api_url", "https://api-staging.pingen.com")
	v.SetDefault("pingen.scope", "letter")
	v.SetDefault("pingen.product", "fast")
	v.SetDefault("pingen.address_position", "left")
	v.SetDefault("pingen.print_mode", "simplex")
	v.SetDefault("pingen.print_spectrum", "color")
}

// GetAddr returns the server address in host:port format
func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
