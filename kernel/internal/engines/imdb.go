package engines

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
)

const (
	paging         = false
	suggestion_url = "https://v2.sg.media-imdb.com/suggestion/%s/%s.json"
	href_base      = "https://imdb.com/%s/%s"

	EngineNameIMDB = "imdb"
)

var (
	categories        = []string{"movies"}
	search_categories = map[string]string{"nm": "name", "tt": "title", "kw": "keyword", "co": "company", "ep": "episode"}
)

type imdb struct {
}

func init() {
	search.RegisterEngine(EngineNameIMDB, &imdb{}, CategoryGeneral)
}

func (e *imdb) Request(ctx context.Context, opts *search.Options) error {
	log := slog.With("func", "imdb.Request")

	// imdb engine does not support finding the next page.
	// So no result is returned when the number of pages requested is greater than 1.
	if opts.PageNo > 1 {
		return nil
	}

	query := strings.ToLower(strings.Replace(opts.Query, " ", "_", -1))
	if len(query) == 0 {
		return nil
	}

	letter := query[0:1]
	if !utf8.ValidString(letter) {
		letter = "x"
	}
	opts.Url = fmt.Sprintf(suggestion_url, letter, query)
	log.DebugContext(ctx, "request", "url", opts.Url)
	return nil
}

func (e *imdb) Response(ctx context.Context, opts *search.Options, resp []byte) (*result.Result, error) {
	log := slog.With("func", "imdb.Response")

	log.DebugContext(ctx, "response", "resp", string(resp))
	m, err := objx.FromJSON(string(resp))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}
	res := result.CreateResult(EngineNameIMDB, opts.PageNo)
	m.Get("d").EachObjxMap(func(i int, v objx.Map) bool {
		entry := v.Get("id").Str()
		if len(entry) < 2 {
			return true
		}
		category := entry[0:2]
		if _, ok := search_categories[category]; !ok {
			return true
		}

		title := v.Get("l").Str()
		if v.Has("q") {
			title += fmt.Sprintf("(%s)", v.Get("q").Str())
		}

		content := ""
		if v.Has("rank") {
			content += fmt.Sprintf("(%v)", v.Get("rank").Float64())
		}
		if v.Has("y") {
			content += fmt.Sprintf("%v - ", v.Get("y").Float64())
		}
		if v.Has("s") {
			content += v.Get("s").Str()
		}

		image := v.Get("i").ObjxMap().Get("imageUrl").Str()
		if image != "" {
			// get thumbnail
			image = strings.Replace(image, "._V1_.", "._V1_QL75_UX280_CR0,0,280,414_.", -1)
		}

		res.AppendData(result.Data{
			Engine:  EngineNameIMDB,
			Title:   title,
			Url:     fmt.Sprintf(href_base, search_categories[category], entry),
			Content: content,
			ImgSrc:  image,
			Query:   opts.Query,
		})
		return true
	})

	return res, nil
}
