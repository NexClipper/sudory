package config

import "github.com/jinzhu/configor"

type Config struct {
	APPName string

	Host struct {
		Port int32
	}

	Database struct {
		Type            string
		DSN             string
		MaxOpenConns    int
		MaxIdleConns    int
		MaxConnLifeTime int
		ShowSQL         bool
		LogLevel        string
	}
}

func New(configPath string) (*Config, error) {
	c := &Config{}
	if err := configor.Load(c, configPath); err != nil {
		return nil, err
	}
	return c, nil
}
