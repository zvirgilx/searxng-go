package network

import (
	"net/http"
	"net/url"

	"github.com/zvirgilx/searxng-go/kernel/config"
	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

var client *http.Client

func InitClient(conf config.NetworkConfig) error {
	client = httputil.NewClient(conf.Timeout)

	if conf.ProxyUrl != "" {
		if err := setProxy(conf.ProxyUrl); err != nil {
			return err
		}
	}

	return nil
}

func GetClient() *http.Client {
	return client
}

func setProxy(proxy string) error {
	parsedU, err := url.Parse(proxy)
	if err != nil {
		return err
	}

	t, ok := client.Transport.(*http.Transport)
	if client.Transport != nil && ok {
		t.Proxy = http.ProxyURL(parsedU)
		t.DisableKeepAlives = true
	} else {
		client.Transport = &http.Transport{
			Proxy:             http.ProxyURL(parsedU),
			DisableKeepAlives: true,
		}
	}
	return nil
}
