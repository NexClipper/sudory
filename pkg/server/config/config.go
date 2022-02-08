package config

import "github.com/jinzhu/configor"

type Config struct {
	APPName string

	Host struct {
		Port int32
	}

	Database struct {
		Type            string
		Protocol        string `default:"tcp"`
		Host            string `env:"SUDORY_DB_HOST"`
		Port            string `env:"SUDORY_DB_PORT"`
		DBName          string `env:"SUDORY_DB_SCHEME"`
		Username        string `env:"SUDORY_DB_SERVER_USERNAME"`
		Password        string `env:"SUDORY_DB_SERVER_PASSWORD"`
		MaxOpenConns    int
		MaxIdleConns    int
		MaxConnLifeTime int
		ShowSQL         bool
		LogLevel        string
	}
}

func New(c *Config, configPath string) (*Config, error) {
	if c == nil {
		c = &Config{}
	}
	if err := configor.Load(c, configPath); err != nil {
		return nil, err
	}
	return c, nil
}
