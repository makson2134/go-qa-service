package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Port         string        `env:"PORT" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string        `env:"POSTGRES_HOST" env-required:"true"`
	Port            string        `env:"POSTGRES_PORT" env-required:"true"`
	User            string        `env:"POSTGRES_USER" env-required:"true"`
	Database        string        `env:"POSTGRES_DB" env-required:"true"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func (d *DatabaseConfig) GetDSN() (string, error) {
	passwordFile := "/run/secrets/db-password" // #nosec G101
	data, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password from secrets: %w", err)
	}

	password := strings.TrimSpace(string(data))

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, password, d.Database)

	return dsn, nil
}

func Load(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
