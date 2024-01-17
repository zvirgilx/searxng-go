package engines

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zvirgilx/searxng-go/kernel/internal/complete"
	"github.com/zvirgilx/searxng-go/kernel/internal/engines/traits"
	"github.com/zvirgilx/searxng-go/kernel/internal/locale"
	"github.com/zvirgilx/searxng-go/kernel/internal/network"
	"github.com/zvirgilx/searxng-go/kernel/internal/result"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
	"github.com/zvirgilx/searxng-go/kernel/internal/util"
	httputil "github.com/zvirgilx/searxng-go/kernel/internal/util/http"
)

const (
	EngineNameGoogle = "google"
)

var (
	googleTimeRangeMap = map[string]string{
		"day":   "d",
		"week":  "w",
		"month": "m",
		"year":  "y"}
)

type google struct{}

func init() {
	search.RegisterEngine(EngineNameGoogle, &google{}, CategoryGeneral)
	complete.RegisterCompleter(EngineNameGoogle, &google{})
}

func (g *google) Request(ctx context.Context, opts *search.Options) error {
	log := slog.With("func", "google.Request")

	queryParams := url.Values{}
	queryParams.Set("q", opts.Query)
	queryParams.Set("filter", "0")
	queryParams.Set("start", strconv.Itoa((opts.PageNo-1)*10))
	queryParams.Set("async", "use_ac:true,_fmt:prog")

	info := GetGoogleInfo(map[string]string{"locale": opts.Locale})
	param := info["param"].(map[string]string)
	queryParams.Set("hl", param["hl"])
	queryParams.Set("lr", param["lr"])
	queryParams.Set("cr", param["cr"])

	queryUrl := fmt.Sprintf("https://%s/search?%s", info["subdomain"], queryParams.Encode())

	if t, ok := googleTimeRangeMap[opts.TimeRange]; ok {
		qP := url.Values{}
		qP.Set("tbs", "qdr:"+t)
		queryUrl += "&" + qP.Encode()
	}

	opts.Url = queryUrl
	log.DebugContext(ctx, "request", "url", opts.Url)
	return nil

}

func (g *google) Response(ctx context.Context, opts *search.Options, bytes []byte) (*result.Result, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bytes)))
	if err != nil {
		return nil, errors.New("error parsing document")
	}

	res := result.CreateResult(EngineNameGoogle, opts.PageNo)
	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").First().Text()
		link, _ := s.Find("a").First().Attr("href")
		content := s.Find(".VwiC3b").First().Text()

		// ignore empty title result
		if title == "" {
			return
		}

		// ignore redirect url
		if strings.HasPrefix(link, "/search?") {
			return
		}

		// ignore empty content
		if content == "" {
			return
		}

		res.AppendData(result.Data{
			Engine:  EngineNameGoogle,
			Title:   title,
			Url:     link,
			Content: content,
			Query:   opts.Query,
		})
	})

	doc.Find("div.s75CSd").Each(func(i int, s *goquery.Selection) {
		sug := s.First().Text()
		util.SetAdd(res.Suggestions, sug)
	})

	return res, nil
}

func (g *google) Complete(ctx context.Context, q string, locale string) []complete.Result {
	log := slog.With("func", "google.Complete")

	queryParams := url.Values{}
	queryParams.Set("q", q)
	queryParams.Set("client", "chrome")

	info := GetGoogleInfo(map[string]string{"locale": locale})
	param := info["param"].(map[string]string)
	queryParams.Set("hl", param["hl"])

	queryUrl := fmt.Sprintf("https://%s/complete/search?%s", info["subdomain"], queryParams.Encode())
	resp, err := httputil.Get(ctx, network.GetClient(), queryUrl)
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil
	}
	var data []interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		log.ErrorContext(ctx, "err", err)
		return nil
	}

	if len(data) < 2 {
		log.ErrorContext(ctx, "err", "resp too short")
		return nil
	}
	var results []complete.Result
	for _, v := range data[1].([]interface{}) {
		results = append(results, complete.Result{
			Type: complete.TypeText,
			Text: v.(string),
			Info: "",
		})
	}
	return results
}

func GetGoogleInfo(params map[string]string) map[string]interface{} {
	info := make(map[string]interface{})
	param := make(map[string]string)

	searchLocale := params["locale"]
	if searchLocale == "" {
		searchLocale = traits.LocaleAll
	}

	trait := traits.GetTrait(EngineNameGoogle)

	info["language"] = locale.GetLanguageFromTrait(searchLocale, trait, "lang_en")
	info["country"] = trait.GetRegion(searchLocale)
	info["subdomain"] = "www.google.com"
	if subDomain := trait.GetCustom("supported_domains"); subDomain != nil {
		if d, ok := subDomain[strings.ToUpper(trait.GetRegion(searchLocale))]; ok {
			info["subdomain"] = d
		}
	}

	// The hl (host language) parameter specifies the interface language of the user interface.
	langParts := strings.Split(info["language"].(string), "_")
	param["hl"] = fmt.Sprintf("%s-%s", langParts[len(langParts)-1], info["country"].(string))

	// The lr (language restrict) parameter restricts search results to documents written in a particular language.
	params["lr"] = info["language"].(string)
	if searchLocale == traits.LocaleAll {
		params["lr"] = ""
	}

	// The cr parameter restricts search results to documents originating in a particular country.
	params["cr"] = ""
	if country := info["country"].(string); country != "" {
		params["cr"] = "country" + country
	}

	info["param"] = param
	return info
}
