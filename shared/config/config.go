package config

// AppConfig holds the application-wide configuration settings.
type AppConfig struct {
	Postgres    SQLConfig      `json:"postgres"`
	Clickhouse  SQLConfig      `json:"clickhouse"`
	Twitch      TwitchConfig   `json:"twitch"`
	RabbitMQ    RabbitMQConfig `json:"rabbitmq"`
	Redis       RedisConfig    `json:"redis"`
	Webhook     WebhookConfig  `json:"webhook"`
	Server      ServerConfig   `json:"server"`
	EventServer EventServer    `json:"event_server"`
	Minio       MinioBucketConfig
}

// SQLConfig represents configuration for connecting to a SQL database.
type SQLConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Timeout    int    `json:"timeout"`
	MaxRetries int    `json:"max_retries"`
}

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

// RabbitMQConfig holds the configuration for RabbitMQ connection.
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

// RedisConfig holds the configuration for Redis connection.
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// WebhookConfig holds the configuration for webhook.
type WebhookConfig struct {
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// APIConfig holds the configuration for the API server.
type ServerConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	JwtSecret []byte `json:"secret_key,omitempty"`
	JwtExpiry int    `json:"jwt_expiry"`
}

// EventServer holds the configuration for the event server.
type EventServer struct {
	Host string
	Port string
}

// MinioBucketConfig holds the configuration for Minio bucket.
type MinioBucketConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
	BucketName string `json:"BucketName"`
	MaxFileSize int64  `json:"MaxFileSize"`
}
