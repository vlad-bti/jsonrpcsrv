package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		JsonRpc `yaml:"jsonrpc"`
		Log     `yaml:"logger"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// JsonRpc -.
	JsonRpc struct {
		Port  string `env-required:"true" yaml:"port"  env:"JSON_RPC_PORT"`
		Route string `env-required:"true" yaml:"route" env:"JSON_RPC_ROUTE"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// Create config structure
	cfg := &Config{}

	// Open config file
	file, err := os.Open("./config/config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
