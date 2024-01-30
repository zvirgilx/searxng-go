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
	baseUrl    = "127.0.0.1:9200"
	index      = ""
	searchUrl  = ""
	queryType  = "match"
	queryField = "title"
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
	if field := conf.QueryField; field != "" {
		queryField = field
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
	opts.Body = f(queryField, opts.Query)

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
		content := r.Get("content").Str()
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

type queryFunc func(key, value string) string

var availableQueryFunc = map[string]queryFunc{
	"match": matchQuery,
	"term":  termQuery,
}

func matchQuery(key, value string) string {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				key: map[string]interface{}{
					"query": value,
				},
			},
		},
	}

	d, _ := json.Marshal(q)
	return string(d)
}

func termQuery(key, value string) string {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				key: value,
			},
		},
	}
	d, _ := json.Marshal(q)
	return string(d)
}

func (e *elasticSearch) GetName() string {
	return EngineNameElasticSearch
}
