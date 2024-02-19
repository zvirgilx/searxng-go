package result

import (
	"sort"

	"github.com/zvirgilx/searxng-go/kernel/internal/util"
)

type Config struct {
	Score  Score                     `mapstructure:"score"`
	Limits map[string]map[string]int `mapstructure:"limits"`
}

// Result of search
type Result struct {
	MergedData  []*Data   `json:"merged_data"` // MergedData store result from different search engines.
	Suggestions *util.Set `json:"suggestions"` // Suggestions store suggestion from different search engines.
	InfoBox     *InfoBox  `json:"infoBox"`     // InfoBox store information from wikipedia of query

	From   string `json:"-"` // From means the engine name of the search results.
	PageNo int    `json:"-"` // PageNo means the page number of result. PageNo = 1 means first page.
}

// InfoBox of search query from wikipedia(temporary)
type InfoBox struct {
	Title   string              `json:"title"`
	Content string              `json:"content"`
	ImgSrc  string              `json:"img_src"`
	Url     string              `json:"url"`
	UrlList []map[string]string `json:"url_list"`
}

var conf Config

func InitConfig(c Config) {
	conf = c

	loadRule()
}

func CreateResult(from string, page int) *Result {
	return &Result{
		Suggestions: util.NewSet(),
		From:        from,
		PageNo:      page,
	}
}

// Merge engine search result
func (r *Result) Merge(result *Result) {
	result.sortData()

	page := ""
	if r.isFirstPage() {
		page = "first"
	}

	limit := len(result.MergedData)
	if maxSize, ok := conf.Limits[page]; ok {
		if m, have := maxSize[result.From]; have && m < limit {
			limit = m
		}
	}

	r.MergedData = append(r.MergedData, result.MergedData[:limit]...)

	util.SetMerge[string](r.Suggestions, result.Suggestions)

	if r.InfoBox == nil && result.InfoBox != nil {
		r.InfoBox = result.InfoBox
	}

}

func (r *Result) isFirstPage() bool {
	return r.PageNo == 1
}

func (r *Result) AppendData(d *Data) {
	r.MergedData = append(r.MergedData, d.unstructured().doScore())
}

func (r *Result) GetDataSize() int {
	if r == nil {
		return 0
	}
	return len(r.MergedData)
}

func (r *Result) GetSortedData() []*Data {
	r.sortData()
	return r.MergedData
}

func (r *Result) sortData() {
	sort.Slice(r.MergedData, func(i, j int) bool {
		return r.MergedData[i].score > r.MergedData[j].score
	})
}
