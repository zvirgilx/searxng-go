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

	"github.com/PuerkitoBio/goquery"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
)

const (
	EngineNameBingVideos = "bing_videos"
	bingVideosBaseUrl    = "https://www.bing.com/videos/asyncv2"
)

var (
	bingTimeMap = map[string]int{
		"day":   60 * 24,
		"week":  60 * 24 * 7,
		"month": 60 * 24 * 31,
		"year":  60 * 24 * 365,
	}
)

type bingVideo struct{}

func init() {
	search.RegisterEngine(EngineNameBingVideos, &bingVideo{}, CategoryGeneral)
	search.RegisterEngine(EngineNameBingVideos, &bingVideo{}, CategoryVideo)
}

func (e *bingVideo) Request(ctx context.Context, opts *search.Options) error {
	log := slog.With("func", "bing_videos.Request")

	// example: https://www.bing.com/videos/asyncv2?q=test&async=content&first=1&count=35
	queryParams := url.Values{}
	queryParams.Set("q", opts.Query)
	queryParams.Set("async", "content")
	queryParams.Set("first", strconv.Itoa((opts.PageNo-1)*10+1))
	queryParams.Set("count", "10")

	// example: one day (60 * 24 minutes) '&qft= filterui:videoage-lt1440&form=VRFLTR'

	timeRange := opts.TimeRange
	if timeRange != "" {
		queryParams.Set("form", "VRFLTR")
		queryParams.Set("qft", fmt.Sprintf(" filterui:videoage-lt%v", bingTimeMap[timeRange]))
	}

	opts.Url = bingVideosBaseUrl + "?" + queryParams.Encode()
	log.DebugContext(ctx, "request", "url", opts.Url)
	return nil
}

var bingVideoRespRegex = regexp.MustCompile(`(?s)<div class="vidres".*`)

func (e *bingVideo) Response(ctx context.Context, opts *search.Options, resp []byte) (*result.Result, error) {
	log := slog.With("func", "bing_videos.Response")

	matches := bingVideoRespRegex.FindStringSubmatch(string(resp))
	if len(matches) == 0 {
		return nil, errors.New("error parsing document")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(matches[0]))
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil, err
	}

	res := result.CreateResult(EngineNameBingVideos, opts.PageNo)
	doc.Find("div[class=dg_u] div[id^='mc_vtvc_video']").Each(func(i int, s *goquery.Selection) {
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

		res.AppendData(result.Data{
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
