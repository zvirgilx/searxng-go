package network

import (
	"net/http"
	"net/url"
	"time"
)

// Client is an encapsulation and extension of the http client.
type Client struct {
	Client *http.Client
}

type Config struct {
	Timeout  time.Duration `mapstructure:"timeout"`
	ProxyUrl string        `mapstructure:"proxy_url"`
}

// DefaultClient return the default http client.
func DefaultClient() *Client {
	return NewClient(nil)
}

// NewClient create a client base on network config.
// if config is nil, it will return a default client.
func NewClient(config *Config) *Client {
	if config == nil {
		return &Client{http.DefaultClient}
	}

	transport := http.DefaultTransport
	if config.ProxyUrl != "" {
		if parsedU, err := url.Parse(config.ProxyUrl); err == nil {
			transport = &http.Transport{
				Proxy:             http.ProxyURL(parsedU),
				DisableKeepAlives: true,
			}
		}
	}

	if transport != http.DefaultTransport || config.Timeout > 0 {
		return &Client{&http.Client{Timeout: config.Timeout, Transport: transport}}
	}

	return &Client{http.DefaultClient}
}

func (c *Client) Get() *Request {
	return NewRequest(c).Method(http.MethodGet)
}
