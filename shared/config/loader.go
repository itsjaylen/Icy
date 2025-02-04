package config

import (
	"encoding/json"
	"fmt"
	"itsjaylen/IcyConfig/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

var k = koanf.New(".")

// LoadConfig is the main function to load the application configuration.
func LoadConfig(configType string) (*AppConfig, error) {
	switch configType {
	case "debug":
		fmt.Println("Loading debug configuration")
		return loadDebugConfig()
	case "release":
		fmt.Println("Release mode: Configuration should be loaded from environment variables")
		return loadReleaseConfig()
	default:
		return nil, fmt.Errorf("unknown configuration type: %s", configType)
	}
}

// loadDebugConfig handles loading or creating the debug configuration file.
func loadDebugConfig() (*AppConfig, error) {
	configDir := "config"
	configFile := filepath.Join(configDir, "debug.json")

	if err := ensureConfigDirExists(configDir); err != nil {
		return nil, err
	}

	if err := ensureConfigFileExists(configFile); err != nil {
		return nil, err
	}

	return readConfigFromFile(configFile)
}

// ensureConfigDirExists ensures that the configuration directory exists.
func ensureConfigDirExists(configDir string) error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}
	}
	return nil
}

// ensureConfigFileExists ensures that the debug configuration file exists, creating it if necessary.
func ensureConfigFileExists(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			return fmt.Errorf("error creating config file: %w", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(DefaultConfig); err != nil {
			return fmt.Errorf("error writing default config: %w", err)
		}
		fmt.Printf("Created default config file: %s\n", configFile)
	}
	return nil
}

// readConfigFromFile reads and decodes the configuration from a file.
func readConfigFromFile(configFile string) (*AppConfig, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return &config, nil
}

// loadReleaseConfig handles loading configuration from environment variables and flags.
func loadReleaseConfig() (*AppConfig, error) {
	// Define flags
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	models.PostgresFlags(fs)
	models.ClickhouseFlags(fs)
	models.TwitchFlags(fs)
	models.RabbitMQFlags(fs)
	models.RedisFlags(fs)
	models.WebhookFlags(fs)
	models.ServerFlags(fs)

	// Parse the command line flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	// Load flags into Koanf
	if err := k.Load(posflag.Provider(fs, ".", k), nil); err != nil {
		return nil, fmt.Errorf("error loading flags into koanf: %w", err)
	}

	// Load environment variables into Koanf
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "APP_")), "_", ".")
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading environment variables into koanf: %w", err)
	}

	// Unmarshal the final configuration into the AppConfig struct
	var config AppConfig
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling configuration: %w", err)
	}

	return &config, nil
}
