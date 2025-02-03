package models

import "github.com/spf13/pflag"

// TwitchConfig holds the configuration for Twitch API and chat.
type TwitchConfig struct {
	ClientID         string           `json:"client_id"`
	ClientSecret     string           `json:"client_secret"`
	OauthURI         string           `json:"oauth_uri"`
	TwitchChatConfig TwitchChatConfig `json:"twitch_chat"`
}

// TwitchChatConfig configures Twitch chat settings.
type TwitchChatConfig struct {
	Enabled  bool `json:"enabled"`
	Loopback bool `json:"loopback"`
}

// TwitchFlags defines the command line flags for Twitch configuration.
func TwitchFlags(fs *pflag.FlagSet) {
	fs.String("twitch.client_id", "", "Twitch client ID")
	fs.String("twitch.client_secret", "", "Twitch client secret")
	fs.String("twitch.oauth_uri", "", "Twitch OAuth URI")
	fs.Bool("twitch.chat.enabled", false, "Enable Twitch chat")
	fs.Bool("twitch.chat.loopback", false, "Enable Twitch chat loopback")
}
