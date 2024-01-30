package engine

import httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"

// Options for search.
type Options struct {
	Query     string
	Url       string
	PageNo    int
	TimeRange string
	Locale    string
	Category  string

	Body        string
	HTTPOptions []httputil.RequestOption
}

func (o *Options) SetHTTPOptions(opts ...httputil.RequestOption) {
	if len(opts) == 0 {
		return
	}
	o.HTTPOptions = opts
}
