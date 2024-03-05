package engines

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/url"

	"github.com/mitchellh/mapstructure"
	"github.com/zvirgilx/searxng-go/kernel/internal/network"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const (
	EngineNameElasticSearch = "elastic_search"
)

type elasticSearch struct {
	client *network.Client

	baseUrl     string
	index       string
	queryType   string
	queryFields []string

	maxLengthOfContent int
}

type ElasticSearchConfig struct {
	Enable             bool     `mapstructure:"enable"`
	BaseUrl            string   `mapstructure:"base_url"`              // BaseUrl is elastic search access url.
	Index              string   `mapstructure:"index"`                 // Index used by search.
	QueryType          string   `mapstructure:"query_type"`            // The type of query, such as match,term, etc.
	QueryFields        []string `mapstructure:"query_fields"`          // The fields of the query, such as title, content, etc.
	MaxLengthOfContent int      `mapstructure:"max_length_of_content"` // The maximum length of the content.
}

func init() {
	engine.RegisterGlobalEngine(&elasticSearch{client: network.DefaultClient()}, engine.CategoryGeneral)
}

func (e *elasticSearch) Request(ctx context.Context, opts *engine.Options) error {
	if opts.PageNo > 1 {
		return nil
	}

	base, err := url.Parse(e.baseUrl)
	if err != nil {
		return err
	}

	r := e.client.Get().Base(base).Path(e.index+"/_search").
		Body(getQueryFn(e.queryType)(opts.Query, e.queryFields...)).
		Header("Content-Type", "application/json")

	opts.Request = r
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
		if len(content) > e.maxLengthOfContent {
			content = content[:e.maxLengthOfContent]
		}
		imgSrc := r.Get("poster").Str()
		url := r.Get("url").Str()

		res.AppendData(&result.Data{
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

func (e *elasticSearch) GetName() string {
	return EngineNameElasticSearch
}

func (e *elasticSearch) ApplyConfig(conf engine.Config) error {
	e.client = network.NewClient(conf.Client)

	var esConf *ElasticSearchConfig
	if err := mapstructure.Decode(conf.Extra, &esConf); err != nil {
		return err
	}

	e.baseUrl = esConf.BaseUrl
	e.index = esConf.Index
	e.queryType = esConf.QueryType
	e.queryFields = esConf.QueryFields

	e.maxLengthOfContent = esConf.MaxLengthOfContent
	if e.maxLengthOfContent == 0 {
		e.maxLengthOfContent = 500
	}
	return nil

}

type queryFunc func(value string, keys ...string) []byte

var availableQueryFunc = map[string]queryFunc{
	"match":       matchQuery,
	"multi_match": multiMatchQuery,
}

// matchQuery format match query.
//
//	{
//	   "query": {
//	       "match": {
//	           "title": {
//	               "query": "spider"
//	           }
//	       }
//	   }
//	}
func matchQuery(value string, keys ...string) []byte {
	if len(keys) == 0 {
		return nil
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
	return d
}

// multiMatchQuery format multi match query.
//
//	{
//	   "query": {
//	       "multi_match": {
//	           "fields": [
//	               "title",
//	               "description"
//	           ],
//	           "query": "searxng-go"
//	       }
//	   }
//	}
func multiMatchQuery(value string, keys ...string) []byte {
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  value,
				"fields": keys,
			},
		},
	}
	d, _ := json.Marshal(q)
	return d
}

// It can choose different query types.
// This determines the accuracy of the query results.
func getQueryFn(queryType string) queryFunc {
	f := availableQueryFunc[queryType]
	if f == nil {
		slog.Warn("unknown elasticsearch query function type", slog.String("type", queryType))
		return availableQueryFunc["match"]
	}
	return f
}
