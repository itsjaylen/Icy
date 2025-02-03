package config

var DefaultConfig = AppConfig{
	Postgres: SQLConfig{
		Host:       "localhost",
		Port:       "5432",
		User:       "postgres",
		Password:   "password",
		Database:   "exampledb",
		Timeout:    30,
		MaxRetries: 5,
	},
	Clickhouse: SQLConfig{
		Host:       "localhost",
		Port:       "9000",
		User:       "default",
		Password:   "password",
		Database:   "exampledb",
		Timeout:    30,
		MaxRetries: 5,
	},
	Twitch: TwitchConfig{
		ClientID:     "",
		ClientSecret: "",
		OauthURI:     "",
		TwitchChatConfig: TwitchChatConfig{
			Enabled:  false,
			Loopback: false,
		},
	},
	RabbitMQ: RabbitMQConfig{
		Host:     "localhost",
		Port:     "5672",
		User:     "guest",
		Password: "guest",
		Enabled:  false,
	},
	Redis: RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
	},
	Webhook: WebhookConfig{
		URL:     "",
		Enabled: false,
	},
	Server: ServerConfig{
		Host:      "localhost",
		Port:      "8080",
		SecretKey: "your-secret-key",
	},
}