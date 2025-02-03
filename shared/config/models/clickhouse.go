package models

import "github.com/spf13/pflag"

// ClickhouseFlags defines the command line flags for Clickhouse configuration.
func ClickhouseFlags(fs *pflag.FlagSet) {
	fs.String("clickhouse.host", "", "Clickhouse host")
	fs.String("clickhouse.port", "", "Clickhouse port")
	fs.String("clickhouse.user", "", "Clickhouse user")
	fs.String("clickhouse.password", "", "Clickhouse password")
	fs.String("clickhouse.database", "", "Clickhouse database")
	fs.Int("clickhouse.timeout", 0, "Clickhouse connection timeout")
	fs.Int("clickhouse.max_retries", 0, "Clickhouse max retries")
}
