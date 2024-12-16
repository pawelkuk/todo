package config

import "time"

type Config struct {
	DBPath         string        `env:"DB_PATH"`
	Editor         string        `env:"EDITOR"`
	SessionRefresh time.Duration `env:"SESSION_REFRESH"`
}
