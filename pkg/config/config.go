package config

type Config struct {
	DBPath string `env:"DB_PATH"`
	Editor string `env:"EDITOR"`
}
