package config

import (
	"time"

	"github.com/jinzhu/configor"
)

type Config struct {
	APPName string `default:"sudory-server"`

	Host struct {
		Port                   int32  `env:"SUDORY_HOST_PORT"             yaml:"port"             default:"8099"`
		TlsEnable              bool   `env:"SUDORY_HOST_TLS_ENABLE"       yaml:"tls-enable"       default:"false"`
		TlsCertificateFilename string `env:"SUDORY_HOST_TLS_CRT_FILENAME" yaml:"tls-crt-filename" default:"server.crt"`
		TlsPrivateKeyFilename  string `env:"SUDORY_HOST_TLS_KEY_FILENAME" yaml:"tls-key-filename" default:"server.key"`

		XAuthToken bool `default:"false"`
	}

	Database struct {
		Type            string `default:"mysql"`
		Protocol        string `default:"tcp"`
		Host            string `env:"SUDORY_DB_HOST"`
		Port            string `env:"SUDORY_DB_PORT"`
		DBName          string `env:"SUDORY_DB_SCHEME"`
		Username        string `env:"SUDORY_DB_SERVER_USERNAME"`
		Password        string `env:"SUDORY_DB_SERVER_PASSWORD"`
		MaxOpenConns    int    `default:"15"`
		MaxIdleConns    int    `default:"5"`
		MaxConnLifeTime int    `default:"1"`
		ShowSQL         bool   `default:"false"`
		LogLevel        string `default:"warn"`
	}

	Migrate struct {
		Source string `yaml:"source" default:"./schema"`
	} `yaml:"migrate"`

	CORSConfig struct {
		AllowOrigins string `env:"SUDORY_CORSCONFIG_ALLOW_ORIGINS" yaml:"allow-origins,omitempty"`
		AllowMethods string `env:"SUDORY_CORSCONFIG_ALLOW_METHODS" yaml:"allow-methods,omitempty"`
	} `yaml:"cors-config,omitempty"`

	Encryption string `yaml:"encryption" default:"enigma.yml"`

	// @Deprecated
	Events string `yaml:"events" default:"events.yml"`

	// @Deprecated
	RespitePeriod time.Duration `default:"60m" yaml:"respite-period"` //(minute) 0: no use
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
