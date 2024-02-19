package engines

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/zvirgilx/searxng-go/kernel/internal/network"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/engines/traits"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const EngineNameWikipedia = "wikipedia"

type wikipedia struct {
	client *network.Client
}

func init() {
	engine.RegisterGlobalEngine(&wikipedia{client: network.DefaultClient()}, engine.CategoryGeneral)
}

func (w *wikipedia) Request(ctx context.Context, opts *engine.Options) error {
	// if not the first page, not request wikipedia for information.
	if opts.PageNo > 1 {
		return nil
	}

	_, wikiNetLoc := getWikiInfo(opts.Locale)
	title := url.QueryEscape(opts.Query)

	base, err := url.ParseRequestURI(wikiNetLoc)
	if err != nil {
		return err
	}

	opts.Request = w.client.Get().Base(base).Path(fmt.Sprintf("api/rest_v1/page/summary/%s", title))

	return nil
}

func (w *wikipedia) Response(ctx context.Context, opts *engine.Options, resp []byte) (*result.Result, error) {
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

	wikipediaLink := m.Get("content_urls").ObjxMap().
		Get("desktop").ObjxMap().
		Get("page").
		Str()
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

func (w *wikipedia) GetName() string {
	return EngineNameWikipedia
}

func (i *wikipedia) ApplyConfig(conf engine.Config) error {
	i.client = network.NewClient(conf.Client)
	return nil
}
