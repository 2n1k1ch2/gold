package config

import (
	"errors"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"time"
)

type Config struct {
	// Address of TCP Receiver
	ReceiverAddr string `yaml:"receiver_addr"`

	// storage of snapshot
	RetentionSnapshots int           `yaml:"retention_snapshots"`
	RetentionWindow    time.Duration `yaml:"retention_window"`

	// Release-tags
	ReleaseEnv []string `yaml:"release_env"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		ReceiverAddr:       "0.0.0.0:8080",
		RetentionSnapshots: 30,
		RetentionWindow:    time.Hour,
		ReleaseEnv:         []string{},
	}
}

func validate(cfg Config) error {
	if cfg.ReceiverAddr == "" {
		return errors.New("receiver_addr is required")
	}

	if cfg.RetentionSnapshots <= 0 {
		return errors.New("retention_snapshots must be > 0")
	}

	if cfg.RetentionWindow < time.Second {
		return errors.New("retention_window is too short")
	}

	return nil
}
