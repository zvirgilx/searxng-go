package engines

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const (
	EngineNameFMovies = "fmovies"
)

type FMovies struct{}

func init() {
	engine.RegisterEngine(EngineNameFMovies, &FMovies{}, engine.CategoryGeneral)
}

func (f *FMovies) Request(ctx context.Context, opts *engine.Options) error {
	if opts.PageNo > 1 {
		return nil
	}
	opts.Url = fmt.Sprintf("https://fmoviesz.to/filter?keyword=%s&sort=most_relevance", opts.Query)
	return nil
}

func (f *FMovies) Response(ctx context.Context, opts *engine.Options, bytes []byte) (*result.Result, error) {
	log := slog.With("func", "fmovies.Response")

	body := string(bytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	res := result.CreateResult(EngineNameFMovies, opts.PageNo)
	doc.Find("body > div.wrapper > main > div > div > aside.main > section > div.movies.items > div.item").
		Each(func(i int, selection *goquery.Selection) {
			redirectU, ok := selection.Find(" div.poster > a").Attr("href")
			if !ok {
				return
			}
			image, ok := selection.Find("img").Attr("data-src")
			if !ok {
				return
			}

			res.AppendData(result.Data{
				Engine: EngineNameFMovies,
				Title:  selection.Find("div.meta > a").Text(),
				Url:    "https://fmoviesz.to/" + redirectU,
				ImgSrc: image,
				Query:  opts.Query,
			})
		})

	return res, nil
}

func (f *FMovies) GetName() string {
	return EngineNameFMovies
}
