package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Cache     CacheConfig     `mapstructure:"Cache"`
	Server    ServerConfig    `mapstructure:"server"`
	App       AppConfig       `mapstructure:"app"`
	ShortCode ShortCodeConfig `mapstructure:"shortcode"`
}

func LoadConfig(filePath string) (*Config, error) {
	viper.SetConfigFile(filePath)
	viper.SetEnvPrefix("URL_SHORTENER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(",", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type DatabaseConfig struct {
	Driver        string `mapstructure:"driver"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	DBName        string `mapstructure:"db_name"`
	SSLMode       string `mapstructure:"ssl_mode"`
	MaxIdleConnes int    `mapstructure:"max_idle_connes"`
	MaxOpenConnes int    `mapstructure:"max_open_connes"`
}

func (d DatabaseConfig) Dsn() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", d.Driver, d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

type CacheConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
}

type AppConfig struct {
	BaseURL         string        `mapstructure:"base_url"`
	DefaultDuration time.Duration `mapstructure:"default_duration"`
	// SyncViewDuration time.Duration `mapstructure:"sync_view_duration"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type ShortCodeConfig struct {
	Length int `mapstructure:"length"`
}
