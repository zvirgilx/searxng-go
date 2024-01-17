package http

import "net/http"

type RequestOption func(o *request)

func (o *request) Request() *http.Request {
	o.request.URL.RawQuery = o.queries.Encode()
	return o.request
}

func WithHeaders(headers map[string]string) RequestOption {
	return func(o *request) {
		for k, v := range headers {
			o.request.Header.Add(k, v)
		}
	}
}
