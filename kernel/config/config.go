package config

import (
	"bytes"
	_ "embed"

	"github.com/spf13/viper"
	"github.com/zvirgilx/searxng-go/kernel/internal/complete"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

//go:embed default.yaml
var defaultConfig []byte

type Config struct {
	Engines  map[string]map[string]engine.Config `mapstructure:"engines"`
	Complete complete.Config                     `mapstructure:"complete"`
	Result   result.Config                       `mapstructure:"result"`
}

var (
	Conf *Config
)

// InitConfig The default configuration will be used first.
// If a custom configuration is specified, changes are merged based on the default configuration.
func InitConfig(path string) error {
	// set default configuration first
	v := viper.New()
	b := bytes.NewReader(defaultConfig)
	v.SetConfigType("yaml")
	if err := v.ReadConfig(b); err != nil {
		return err
	}

	// If a custom configuration file is specified, the configuration file is loaded
	// and different in the custom configuration are overwritten to the default configuration
	if path != "" {
		v.SetConfigFile(path)
		if err := v.MergeInConfig(); err != nil {
			return err
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}
	Conf = cfg
	return nil
}
