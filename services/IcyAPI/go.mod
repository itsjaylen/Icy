module IcyAPI

go 1.23.5

replace itsjaylen/IcyLogger => ../../shared/logger

replace itsjaylen/IcyConfig => ../../shared/config

require itsjaylen/IcyLogger v0.0.0

require (
	github.com/spf13/pflag v1.0.5
	itsjaylen/IcyConfig v0.0.0
)

require github.com/knadh/koanf/v2 v2.1.2 // indirect

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/knadh/koanf/providers/posflag v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	golang.org/x/sys v0.25.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

require (
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/env v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/prometheus/client_golang v1.20.5
)
