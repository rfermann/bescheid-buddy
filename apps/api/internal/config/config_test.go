package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Arrange
	os.Setenv("APP_SERVER_PORT", "8080")
	os.Setenv("APP_SERVER_HOST", "localhost")

	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
	})

	// Act
	cfg, err := config.LoadConfig("")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 30, cfg.Server.ReadTimeout)
	assert.Equal(t, 30, cfg.Server.WriteTimeout)
	assert.Equal(t, 30, cfg.Server.GracefulShutdown)
}

func TestLoadConfig_WithConfigFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "cfg.yaml")
	configContent := `
server:
  port: "8080"
  host: "localhost"
  read_timeout: 60
  write_timeout: 120
  graceful_shutdown: 45
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 60, cfg.Server.ReadTimeout)
	assert.Equal(t, 120, cfg.Server.WriteTimeout)
	assert.Equal(t, 45, cfg.Server.GracefulShutdown)
}

func TestLoadConfig_WithEnvironmentVariables(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	configPath := filepath.Join(tempDir, "cfg.yaml")
	configContent := `
server:
  port: "8080"
  host: "localhost"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_SERVER_HOST", "0.0.0.0")
	os.Setenv("APP_SERVER_READ_TIMEOUT", "45")

	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
		os.Unsetenv("APP_SERVER_READ_TIMEOUT")
	})

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 45, cfg.Server.ReadTimeout)
	assert.Equal(t, 30, cfg.Server.WriteTimeout) // Should remain default
}

func TestLoadConfig_WithEnvironmentSpecificConfig(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create environment-specific cfg
	envConfigContent := `
server:
  port: "8081"
  host: "staging.host"
  read_timeout: 45
`
	err = os.WriteFile(filepath.Join(tempDir, "config.staging.yaml"), []byte(envConfigContent), 0644)
	require.NoError(t, err)

	os.Setenv("APP_ENV", "staging")

	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
	})

	// Act
	cfg, err := config.LoadConfig("")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8081", cfg.Server.Port)
	assert.Equal(t, "staging.host", cfg.Server.Host)
	assert.Equal(t, 45, cfg.Server.ReadTimeout)
}

func TestLoadConfig_ConfigValidation(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "cfg.yaml")

	// Create cfg missing required fields
	configContent := `
server:
  read_timeout: 30
  write_timeout: 30
  graceful_shutdown: 30
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "configuration validation failed")
	assert.Contains(t, err.Error(), "server.port is required")
	assert.Contains(t, err.Error(), "server.host is required")
}

func TestLoadConfig_ConfigFileNotFound(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
	})

	// Act
	cfg, err := config.LoadConfig("nonexistent.yaml")

	// Assert
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadConfig_NoConfigFileWithEnvVars(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	os.Setenv("APP_SERVER_PORT", "8080")
	os.Setenv("APP_SERVER_HOST", "localhost")

	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
	})

	// Act
	cfg, err := config.LoadConfig("")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
}

func TestGetAddr(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	// Act
	addr := cfg.GetAddr()

	// Assert
	assert.Equal(t, "localhost:8080", addr)
}

func TestLoadConfig_UnmarshalError_InvalidStructure(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "cfg.yaml")

	// Create a cfg file with invalid YAML that can be read but not unmarshaled
	// Using invalid YAML structure that doesn't match our Config struct
	configContent := `
server:
  port: "8080"
  host: "localhost"
  # This will cause unmarshaling error - read_timeout should be int, not string
  read_timeout: "not_a_number"
  write_timeout: "also_not_a_number"
  graceful_shutdown: "still_not_a_number"
  # Add an unknown field to test strict unmarshaling
  unknown_field: "this_should_not_be_here"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

func TestLoadConfig_UnmarshalError_InvalidYAML(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a file with syntactically invalid YAML
	configContent := `
server:
  port: "8080"
  host: "localhost"
read_timeout: 30
  invalid_yaml_structure: [
    missing_closing_bracket
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadConfig_EnvironmentVariablesOverrideConfigFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "cfg.yaml")
	configContent := `
server:
  port: "8080"
  host: "localhost"
  read_timeout: 30
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_HOST")
	})

	// Set environment variables to override cfg file
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_SERVER_HOST", "0.0.0.0")

	// Act
	cfg, err := config.LoadConfig(configPath)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 30, cfg.Server.ReadTimeout) // Should remain from cfg file
}
