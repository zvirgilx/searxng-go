package search

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/metrics"
	"github.com/zvirgilx/searxng-go/kernel/internal/network"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	"github.com/zvirgilx/searxng-go/kernel/internal/util"
	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

func Search(ctx context.Context, options engine.Options) *result.Result {
	log := slog.With("func", "search.Search")

	log.InfoContext(ctx, "starting search", "query", options.Query)

	enableEngines := engine.GetEnginesByCategory(options.Category)
	if len(enableEngines) == 0 {
		log.WarnContext(ctx, "engines not found", "category", options.Category)
		return &result.Result{}
	}

	resCh := make(chan *result.Result, len(enableEngines))

	w := &sync.WaitGroup{}
	for _, e := range enableEngines {
		w.Add(1)
		go func(opts engine.Options, e engine.Engine) {
			defer w.Done()
			defer util.RecoverFromPanic()
			r, err := process(ctx, opts, e)
			if err != nil {
				log.ErrorContext(ctx, "process error", slog.String("engine", e.GetName()), slog.String("err", err.Error()))
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

func process(ctx context.Context, options engine.Options, e engine.Engine) (r *result.Result, err error) {
	log := slog.With("func", "search.process")

	start := time.Now()

	defer func() {
		status := "ok"
		if err != nil {
			status = "error"
		}

		metrics.EnginesResponseCounter.WithLabelValues(e.GetName(), status).Observe(time.Since(start).Seconds())
		metrics.EnginesSearchResultCounter.WithLabelValues(e.GetName()).Add(float64(r.GetDataSize()))

	}()

	if err = e.Request(ctx, &options); err != nil {
		log.ErrorContext(ctx, "request error", slog.String("engine", e.GetName()), slog.String("err", err.Error()))
		return nil, err
	}

	if options.Url == "" {
		return nil, nil
	}

	body, err := httputil.Get(ctx, network.GetClient(), options.Url, strings.NewReader(options.Body), options.HTTPOptions...)
	if err != nil {
		log.ErrorContext(ctx, "http get error", slog.String("engine", e.GetName()), slog.String("err", err.Error()))
		return nil, err
	}

	r, err = e.Response(ctx, &options, body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func VerifySearchOptions(c *gin.Context) (engine.Options, error) {
	q, ok := c.GetQuery("q")
	if !ok {
		return engine.Options{}, errors.New("empty query input")
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
			return engine.Options{}, errors.New("page number error")
		}
		pageNum = num
	}

	category, ok := c.GetQuery("category")
	if !ok {
		category = "general"
	}

	return engine.Options{
		Query:    q,
		PageNo:   pageNum,
		Locale:   lang,
		Category: category,
	}, nil
}
