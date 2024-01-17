package network

import (
	"net/http"
	"time"

	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

var client *http.Client

func InitClient(timeout time.Duration) {
	client = httputil.NewClient(timeout)
}

func GetClient() *http.Client {
	return client
}
