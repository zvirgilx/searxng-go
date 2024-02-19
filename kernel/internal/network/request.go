package network

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Request is a chain-style implementation of an http request.
// It provides a convenient way to set the request path, parameters,
// request headers, etc.
type Request struct {
	c *Client

	method  string
	base    *url.URL
	path    string
	params  url.Values
	headers http.Header

	timeout time.Duration

	body []byte
}

func NewRequest(c *Client) *Request {
	var timeout time.Duration
	if c.Client != nil {
		timeout = c.Client.Timeout
	}

	r := &Request{
		c:       c,
		timeout: timeout,
	}
	return r
}

func (r *Request) Method(method string) *Request {
	r.method = method
	return r
}

func (r *Request) Base(base *url.URL) *Request {
	r.base = base
	return r
}

func (r *Request) Path(path string) *Request {
	r.path = path
	return r
}

func (r *Request) Header(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

func (r *Request) Param(key string, value string) *Request {
	if r.params == nil {
		r.params = url.Values{}
	}
	r.params.Del(key)
	r.params.Add(key, value)
	return r
}

func (r *Request) Body(b []byte) *Request {
	r.body = b
	return r
}

// URL returns the url according to the Request.
func (r *Request) URL() *url.URL {
	u := r.base

	query := url.Values{}
	for k, vs := range r.params {
		for _, v := range vs {
			query.Add(k, v)
		}
	}

	u.Path = r.path

	u.RawQuery = query.Encode()
	return u
}

// newHTTPRequest returns a build-in http request from Request.
func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	var body io.Reader
	if r.body != nil {
		body = bytes.NewReader(r.body)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, r.URL().String(), body)
	if err != nil {
		return nil, err
	}
	req.Header = r.headers
	return req, nil
}

// request initiates an http request and obtains the response result.
func (r *Request) request(ctx context.Context, fn func(r *http.Request, resp *http.Response)) error {
	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}

	fn(req, resp)

	return nil

}

type Result struct {
	Body       []byte
	Err        error
	StatusCode int
}

// resultForResponse parse the http response.
// if Result.Err is nil means a successful response got.
func (r *Request) resultForResponse(resp *http.Response) Result {
	var body []byte
	if resp.Body != nil {
		d, err := io.ReadAll(resp.Body)
		if err != nil {
			return Result{
				Err: fmt.Errorf("error happen when reading response Body. error: %w", err),
			}
		}
		body = d
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		err := fmt.Errorf("status code of response is not ok. status code: %d", resp.StatusCode)
		return Result{
			Body:       body,
			Err:        err,
			StatusCode: resp.StatusCode,
		}
	}
	return Result{
		Body:       body,
		StatusCode: resp.StatusCode,
	}
}

// Do execute request.
func (r *Request) Do(ctx context.Context) Result {
	var result Result
	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result = r.resultForResponse(resp)
	})
	if err != nil {
		return Result{Err: err}
	}
	return result
}
