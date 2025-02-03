package models


import "github.com/spf13/pflag"

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

// PostgresFlags defines the command line flags for PostgreSQL configuration.
func PostgresFlags(fs *pflag.FlagSet) {
	fs.String("postgres.host", "", "Postgres host")
	fs.String("postgres.port", "", "Postgres port")
	fs.String("postgres.user", "", "Postgres user")
	fs.String("postgres.password", "", "Postgres password")
	fs.String("postgres.database", "", "Postgres database")
	fs.Int("postgres.timeout", 0, "Postgres connection timeout")
	fs.Int("postgres.max_retries", 0, "Postgres max retries")
}
