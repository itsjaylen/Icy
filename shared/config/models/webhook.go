package models

import "github.com/spf13/pflag"

// WebhookConfig holds the configuration for webhook.
type WebhookConfig struct {
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// WebhookFlags defines the command line flags for Webhook configuration.
func WebhookFlags(fs *pflag.FlagSet) {
	fs.String("webhook.url", "", "Webhook URL")
	fs.Bool("webhook.enabled", false, "Enable Webhook")
}
