package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrInvalidHTTPClient = errors.New("invalid http client")
)

type request struct {
	request *http.Request
	queries url.Values
	maxTry  int
}

// NewClient create new http client
func NewClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          300,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			MaxIdleConnsPerHost:   300,
		},
		Timeout: timeout,
	}
	return client
}

// Get HTTP GET
func Get(ctx context.Context, cli *http.Client, url string, reqBody io.Reader, opts ...RequestOption) ([]byte, error) {
	if cli == nil {
		return nil, ErrInvalidHTTPClient
	}

	log := slog.With("func", "Get")
	o, err := newRequest(ctx, http.MethodGet, url, reqBody, opts...)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	for i := 1; i <= o.maxTry; i++ {
		if resp, err = cli.Do(o.Request()); err == nil {
			break
		}
		if i < o.maxTry {
			log.WarnContext(ctx, "http get failed", slog.String("url", url), slog.Int("i", i), slog.String("err", err.Error()))
			time.Sleep(time.Millisecond * (2 << i))
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("http code: %d", resp.StatusCode)
	}

	return body, nil
}

func newRequest(ctx context.Context, method, url string, body io.Reader, opts ...RequestOption) (*request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	o := &request{
		request: req,
		queries: req.URL.Query(),
		maxTry:  1,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o, nil
}
