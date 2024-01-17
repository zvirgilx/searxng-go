package search

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zvirgilx/searxng-go/kernel/internal/network"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

type Engine interface {
	Request(context.Context, *Options) error
	Response(context.Context, *Options, []byte) (*result.Result, error)
}

var engines = map[string]map[string]Engine{}

// RegisterEngine registers a search engine
func RegisterEngine(name string, engine Engine, category string) {
	if engines[category] == nil {
		engines[category] = map[string]Engine{}
	}
	engines[category][name] = engine
}

// Options for search
type Options struct {
	Query     string
	Url       string
	PageNo    int
	TimeRange string
	Locale    string
	Category  string
}

func Search(ctx context.Context, options Options) *result.Result {
	log := slog.With("func", "search.Search")

	log.InfoContext(ctx, "starting search", "query", options.Query)

	enableEngines := engines[options.Category]
	if len(enableEngines) == 0 {
		log.WarnContext(ctx, "engines not found", "category", options.Category)
		return &result.Result{}
	}

	resCh := make(chan *result.Result, len(enableEngines))

	w := &sync.WaitGroup{}
	for _, e := range enableEngines {
		w.Add(1)
		go func(opts Options, e Engine) {
			defer w.Done()
			r, err := process(ctx, opts, e)
			if err != nil {
				log.ErrorContext(ctx, "err", err)
				return
			}
			resCh <- r
		}(options, e)
	}
	w.Wait()
	close(resCh)

	result := result.CreateResult("", options.PageNo)
	for r := range resCh {
		if r == nil {
			continue
		}
		result.Merge(r)
	}

	return result
}

func process(ctx context.Context, options Options, engine Engine) (*result.Result, error) {
	log := slog.With("func", "search.process")

	err := engine.Request(ctx, &options)
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	if options.Url == "" {
		return nil, nil
	}

	// TODO: fix me
	httpOpts := httputil.WithHeaders(map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"})
	body, err := httputil.Get(ctx, network.GetClient(), options.Url, httpOpts)
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	return engine.Response(ctx, &options, body)
}

func VerifySearchOptions(c *gin.Context) (Options, error) {
	q, ok := c.GetQuery("q")
	if !ok {
		return Options{}, errors.New("query not found")
	}

	lang, ok := c.GetQuery("language")
	if !ok {
		lang = "en-US"
	}

	pageNum := 1
	pageNo, ok := c.GetQuery("page_no")
	if ok {
		num, err := strconv.Atoi(pageNo)
		if err != nil || num < 0 {
			return Options{}, errors.New("page number error")
		}
		pageNum = num
	}

	category, ok := c.GetQuery("category")
	if !ok {
		category = "general"
	}

	return Options{
		Query:    q,
		PageNo:   pageNum,
		Locale:   lang,
		Category: category,
	}, nil
}
