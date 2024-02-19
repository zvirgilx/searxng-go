package engine

import (
	"github.com/zvirgilx/searxng-go/kernel/internal/network"
)

// Options for search.
type Options struct {
	Query     string
	Url       string
	PageNo    int
	TimeRange string
	Locale    string
	Category  string

	Request *network.Request
}

type Config struct {
	Enable bool            `mapstructure:"enable"`
	Client *network.Config `mapstructure:"client"`

	Extra interface{} `mapstructure:"extra"`
}
