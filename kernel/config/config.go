package config

import (
	"bytes"
	_ "embed"
	"time"

	"github.com/spf13/viper"
)

//go:embed default.yaml
var defaultConfig []byte

type Config struct {
	Search   SearchConfig   `mapstructure:"search"`
	Complete CompleteConfig `mapstructure:"complete"`
	Network  NetworkConfig  `mapstructure:"network"`
	Result   ResultConfig   `mapstructure:"result"`
}

type SearchConfig struct {
	Engines Engines `mapstructure:"engines"`
}

type Engines struct {
	ElasticSearch ElasticSearch `mapstructure:"elastic_search"`
}

type ElasticSearch struct {
	Enable      bool     `mapstructure:"enable"`       // Whether to enable the engine.
	BaseUrl     string   `mapstructure:"base_url"`     // BaseUrl is elastic search access url.
	Index       string   `mapstructure:"index"`        // Index used by search.
	QueryType   string   `mapstructure:"query_type"`   // The type of query, such as match,term, etc.
	QueryFields []string `mapstructure:"query_fields"` // The fields of the query, such as title, content, etc.
}

type CompleteConfig struct {
	EnableEngines []string `mapstructure:"enable_engines"`
}

type NetworkConfig struct {
	Timeout  time.Duration `mapstructure:"timeout"`
	ProxyUrl string        `mapstructure:"proxy_url"`
}

type ResultConfig struct {
	Score  Score                     `mapstructure:"score"`
	Limits map[string]map[string]int `mapstructure:"limits"`
}

type Score struct {
	Scorer         string         `mapstructure:"scorer"`
	Weight         map[string]int `mapstructure:"weight"`
	MetadataFields []string       `mapstructure:"metadata_fields"`
	Rules          []Rule         `mapstructure:"rules"`
}

type Rule struct {
	Name       string      `mapstructure:"name"`
	Enable     bool        `mapstructure:"enable"`
	Score      int         `mapstructure:"score"`
	Conditions []Condition `mapstructure:"conditions"`
}

type Condition struct {
	Field    string   `mapstructure:"field"`
	Operator string   `mapstructure:"operator"`
	Values   []string `mapstructure:"values"`
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
