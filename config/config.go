package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Config{viper: v}, nil
}

func MustLoad(configPath string) *Config {
	c, err := Load(configPath)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Config) Parse(v any) error {
	return c.viper.Unmarshal(v)
}

func (c *Config) GetString(key string, defaultValue string) string {
	if c.viper.IsSet(key) {
		return c.viper.GetString(key)
	}
	return defaultValue
}

func (c *Config) GetInt(key string, defaultValue int) int {
	if c.viper.IsSet(key) {
		return c.viper.GetInt(key)
	}
	return defaultValue
}

func (c *Config) GetBool(key string, defaultValue bool) bool {
	if c.viper.IsSet(key) {
		return c.viper.GetBool(key)
	}
	return defaultValue
}

func (c *Config) GetDuration(key string, defaultValue time.Duration) time.Duration {
	if c.viper.IsSet(key) {
		return c.viper.GetDuration(key)
	}
	return defaultValue
}
