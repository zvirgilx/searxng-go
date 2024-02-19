package engines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/zvirgilx/searxng-go/kernel/internal/network"

	"github.com/PuerkitoBio/goquery"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
)

const (
	EngineNameBingVideos = "bing_videos"
)

var (
	bingTimeMap = map[string]int{
		"day":   60 * 24,
		"week":  60 * 24 * 7,
		"month": 60 * 24 * 31,
		"year":  60 * 24 * 365,
	}
)

type bingVideo struct {
	client *network.Client
}

func init() {
	engine.RegisterGlobalEngine(&bingVideo{client: network.DefaultClient()}, engine.CategoryGeneral)
}

func (e *bingVideo) Request(ctx context.Context, opts *engine.Options) error {
	// example: https://www.bing.com/videos/asyncv2?q=test&async=content&first=1&count=35
	base, _ := url.Parse("https://www.bing.com")
	req := e.client.Get().Base(base).Path("videos/asyncv2").
		Param("q", opts.Query).
		Param("async", "content").
		Param("first", strconv.Itoa((opts.PageNo-1)*10)).
		Param("count", "10")

	// example: one day (60 * 24 minutes) '&qft= filterui:videoage-lt1440&form=VRFLTR'
	if opts.TimeRange != "" {
		req.Param("form", "VRFLTR").Param("qft", fmt.Sprintf(" filterui:videoage-lt%v", bingTimeMap[opts.TimeRange]))
	}

	opts.Request = req
	return nil
}

// match available html content.
var bingVideoRespRegex = regexp.MustCompile(`(?s)<div class="dg_u".*`)

// Sometimes the html of the first page does not as same format as others.
// So it is compatible with the parsing of the first page.
var bingVideoFirstHTMLRespRegex = regexp.MustCompile(`(?s)<div class="mc_fgvc_u.*`)

func (e *bingVideo) Response(ctx context.Context, opts *engine.Options, resp []byte) (*result.Result, error) {
	log := slog.With("func", "bing_videos.Response")

	body := string(resp)

	// default xpath selector
	xPath := "div[class=dg_u] div[id^='mc_vtvc_video']"

	matches := bingVideoRespRegex.FindStringSubmatch(body)
	if len(matches) == 0 {
		matches = bingVideoFirstHTMLRespRegex.FindStringSubmatch(body)
		if len(matches) == 0 {
			return nil, errors.New("failed to parse bing videos html")
		}

		// compatible xpath selector
		xPath = "div[id^=mc_vtvc__]"
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(matches[0]))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	res := result.CreateResult(EngineNameBingVideos, opts.PageNo)
	doc.Find(xPath).Each(func(i int, s *goquery.Selection) {
		vrhData, exists := s.Find("div.vrhdata").Attr("vrhm")
		if !exists {
			return
		}

		var metadata map[string]interface{}
		if err = json.Unmarshal([]byte(vrhData), &metadata); err != nil {
			return
		}

		info := strings.TrimSpace(s.Find("div.mc_vtvc_meta_block span").Text())
		content := fmt.Sprintf("%s - %s", metadata["du"], info)
		thumbnail, _ := s.Find("div.mc_vtvc_th img").Attr("src")

		res.AppendData(&result.Data{
			Engine:    EngineNameBingVideos,
			Title:     metadata["vt"].(string),
			Url:       metadata["murl"].(string),
			Thumbnail: thumbnail,
			Content:   content,
			Query:     opts.Query,
		})
	})

	return res, nil
}

func (e *bingVideo) GetName() string {
	return EngineNameBingVideos
}

func (e *bingVideo) ApplyConfig(conf engine.Config) error {
	e.client = network.NewClient(conf.Client)
	return nil
}
