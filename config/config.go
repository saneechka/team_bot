package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
	} `yaml:"bot"`

	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
	} `yaml:"database"`

	TelegramAdmins struct {
		Usernames []string `yaml:"usernames"`
	} `yaml:"admins"`
}

func MustLoadConfig(path string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		config.Bot.Token = token
	}

	return config, nil
}
