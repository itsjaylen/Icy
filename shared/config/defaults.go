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
		Port:      "9050",
		JwtSecret: []byte("your-secret-key"),
		JwtExpiry: 3600,
	},
	EventServer: EventServer{
		Host: "localhost",
		Port: "9051",
	},
	Minio: MinioBucketConfig{
		Host:      "localhost",
		Port:      "9000",
		AccessKey: "admin",
		SecretKey: "supersecretpassword",
	},
}
