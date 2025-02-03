package models

import "github.com/spf13/pflag"

// ServerConfig holds the configuration for the API server.
type ServerConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	SecretKey string `json:"secret_key,omitempty"`
}

// ServerFlags defines the command line flags for API configuration.
func ServerFlags(fs *pflag.FlagSet) {
	fs.String("api.host", "", "API host")
	fs.String("api.port", "", "API port")
}
