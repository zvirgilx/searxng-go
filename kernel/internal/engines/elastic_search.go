package engines

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/config"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

const (
	EngineNameElasticSearch = "elastic_search"
)

type elasticSearch struct{}

var (
	baseUrl     = "127.0.0.1:9200"
	index       = ""
	searchUrl   = ""
	queryType   = "match"
	queryFields = []string{"title"}
)

func init() {
	engine.RegisterEngine(EngineNameElasticSearch, &elasticSearch{}, engine.CategoryGeneral)
}

// InitElasticSearch init elastic search engine.
func InitElasticSearch() {
	conf := config.Conf.Search.Engines.ElasticSearch
	if !conf.Enable {
		engine.DisableEngine(engine.CategoryGeneral, EngineNameElasticSearch)
		return
	}
	if u := conf.BaseUrl; u != "" {
		baseUrl = u
	}
	if i := conf.Index; i != "" {
		index = i
	}

	searchUrl = baseUrl + "/" + index + "/_search"

	if t := conf.QueryType; t != "" {
		queryType = t
	}
	if fields := conf.QueryFields; len(fields) != 0 {
		queryFields = fields
	}

}

func (e *elasticSearch) Request(ctx context.Context, opts *engine.Options) error {
	if opts.PageNo > 1 {
		return nil
	}

	opts.Url = searchUrl

	// It can choose different query types.
	// This determines the accuracy of the query results.
	f := availableQueryFunc[queryType]
	if f == nil {
		return fmt.Errorf("elastic search query type not found, type:%s", queryType)
	}

	// The request body is the query condition of es.
	opts.Body = f(opts.Query, queryFields...)

	httpOpts := httputil.WithHeaders(map[string]string{"Content-Type": "application/json"})
	opts.SetHTTPOptions(httpOpts)

	return nil
}

func (e *elasticSearch) Response(ctx context.Context, opts *engine.Options, resp []byte) (*result.Result, error) {
	log := slog.With("func", "elasticSearch.Response")

	log.DebugContext(ctx, "response", "resp", string(resp))
	m, err := objx.FromJSON(string(resp))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	hits := m.Get("hits").ObjxMap()
	if count := hits.Get("total").ObjxMap().Get("value").Int(); count == 0 {
		return nil, nil
	}

	res := result.CreateResult(EngineNameElasticSearch, opts.PageNo)
	hits.Get("hits").EachObjxMap(func(i int, v objx.Map) bool {
		r := v.Get("_source").ObjxMap()

		title := r.Get("title").Str()
		content := r.Get("description").Str()
		imgSrc := r.Get("poster").Str()
		url := r.Get("url").Str()

		res.AppendData(result.Data{
			Engine:  EngineNameElasticSearch,
			Title:   title,
			Url:     url,
			Content: content,
			ImgSrc:  imgSrc,
			Query:   opts.Query,
		})

		return true
	})

	return res, nil
}

type queryFunc func(value string, keys ...string) string

var availableQueryFunc = map[string]queryFunc{
	"match":       matchQuery,
	"multi_match": multiMatchQuery,
}

func matchQuery(value string, keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				keys[0]: map[string]interface{}{
					"query": value,
				},
			},
		},
	}

	d, _ := json.Marshal(q)
	return string(d)
}

func multiMatchQuery(value string, keys ...string) string {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  value,
				"fields": keys,
			},
		},
	}
	d, _ := json.Marshal(q)
	return string(d)
}

func (e *elasticSearch) GetName() string {
	return EngineNameElasticSearch
}
