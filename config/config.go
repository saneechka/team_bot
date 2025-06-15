package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
	} `yaml:"bot"`

	Database struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
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

// GetDatabaseConnectionString returns PostgreSQL connection string
func (c *Config) GetDatabaseConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
