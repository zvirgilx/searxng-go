package engines

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/zvirgilx/searxng-go/kernel/internal/network"

	"github.com/stretchr/objx"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const (
	paging   = false
	hrefBase = "https://imdb.com/%s/%s"

	EngineNameIMDB = "imdb"
)

var (
	search_categories = map[string]string{"nm": "name", "tt": "title", "kw": "keyword", "co": "company", "ep": "episode"}
)

type imdb struct {
	client *network.Client
}

func init() {
	engine.RegisterGlobalEngine(&imdb{client: network.DefaultClient()}, engine.CategoryGeneral)
}

func (i *imdb) Request(ctx context.Context, opts *engine.Options) error {
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

	base, _ := url.ParseRequestURI("https://v2.sg.media-imdb.com")
	opts.Request = i.client.Get().
		Base(base).Path(fmt.Sprintf("suggestion/%s/%s.json", letter, query))
	return nil
}

func (i *imdb) Response(ctx context.Context, opts *engine.Options, resp []byte) (*result.Result, error) {
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

		res.AppendData(&result.Data{
			Engine:  EngineNameIMDB,
			Title:   title,
			Url:     fmt.Sprintf(hrefBase, search_categories[category], entry),
			Content: content,
			ImgSrc:  image,
			Query:   opts.Query,
		})
		return true
	})

	return res, nil
}

func (i *imdb) GetName() string {
	return EngineNameIMDB
}

func (i *imdb) ApplyConfig(conf engine.Config) error {
	i.client = network.NewClient(conf.Client)
	return nil
}
