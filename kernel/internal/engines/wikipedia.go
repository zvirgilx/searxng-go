package engines

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/internal/engines/traits"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
)

const EngineNameWikipedia = "wikipedia"

type wikipedia struct{}

func init() {
	search.RegisterEngine(EngineNameWikipedia, &wikipedia{}, CategoryGeneral)
}

func (w *wikipedia) Request(ctx context.Context, opts *search.Options) error {
	// if not the first page, not request wikipedia for information.
	if opts.PageNo > 1 {
		return nil
	}

	_, wikiNetLoc := getWikiInfo(opts.Locale)
	title := url.QueryEscape(opts.Query)
	opts.Url = fmt.Sprintf("https://%s/api/rest_v1/page/summary/%s", wikiNetLoc, title)

	return nil
}

func (w *wikipedia) Response(ctx context.Context, opts *search.Options, resp []byte) (*result.Result, error) {
	log := slog.With("func", "wikipedia.Response")
	m, err := objx.FromJSON(string(resp))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	title := m.Get("title").Str()
	if title == "" {
		return nil, nil
	}

	wikipediaLink := m.Get("content_urls").ObjxMap().Get("desktop").
		ObjxMap().Get("page").Str()
	if wikipediaLink == "" {
		return nil, nil
	}

	content := m.Get("extract").Str()
	imgSrc := m.Get("thumbnail").ObjxMap().Get("source").Str()

	res := result.CreateResult(EngineNameWikipedia, opts.PageNo)
	infoBox := &result.InfoBox{
		Title:   title,
		Content: content,
		ImgSrc:  imgSrc,
		Url:     wikipediaLink,
		UrlList: []map[string]string{{
			"title": "Wikipedia",
			"url":   wikipediaLink,
		}},
	}
	res.InfoBox = infoBox
	return res, nil
}

func getWikiInfo(locale string) (string, string) {
	trait := traits.GetTrait(EngineNameWikipedia)

	var engTag, wikiNetLoc string
	if engTag = trait.GetRegion(locale); engTag == "" {
		if engTag = trait.GetLanguage(locale); engTag == "" {
			engTag = "en"
		}
	}

	wikiNetLoc = "en.wikipedia.org"
	if wkNetLoc := trait.GetCustom("wiki_netloc"); wkNetLoc != nil {
		if wkNetLoc[engTag] != "" {
			wikiNetLoc = wkNetLoc[engTag]
		}
	}
	return engTag, wikiNetLoc
}
