package result

import (
	"sort"

	"github.com/zvirgilx/searxng-go/kernel/config"
	"github.com/zvirgilx/searxng-go/kernel/internal/util"
)

// Result of search
type Result struct {
	MergedData  []Data    `json:"merged_data"` // MergedData store result from different search engines.
	Suggestions *util.Set `json:"suggestions"` // Suggestions store suggestion from different search engines.
	InfoBox     *InfoBox  `json:"infoBox"`     // InfoBox store information from wikipedia of query

	From   string `json:"-"` // From means the engine name of the search results.
	PageNo int    `json:"-"` // PageNo means the page number of result. PageNo = 1 means first page.
}

// Data of search result
type Data struct {
	Engine    string `json:"engine"`    // Engine is search engine name, means result source.
	Title     string `json:"title"`     // Title is the search result title.
	Url       string `json:"url"`       // Url link to the third party website.
	Content   string `json:"content"`   // Content is a short description.
	ImgSrc    string `json:"img_src"`   // ImgSrc is an image Url, used for poster.
	Thumbnail string `json:"thumbnail"` // Thumbnail Url for some video result.

	// Query is the query of search.
	Query string `json:"-"`

	// Metadata stores data about fields and values in data.
	// This field is not shown externally and is only used internally.
	Metadata map[string]string `json:"-"`

	// Score are scored by enabled scorer on the data.
	Score int `json:"-"`
}

// InfoBox of search query from wikipedia(temporary)
type InfoBox struct {
	Title   string              `json:"title"`
	Content string              `json:"content"`
	ImgSrc  string              `json:"img_src"`
	Url     string              `json:"url"`
	UrlList []map[string]string `json:"url_list"`
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
	if maxSize, ok := config.Conf.Result.Limits[page]; ok {
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

func (r *Result) AppendData(data Data) {
	data.Metadata = buildMetadata(data)
	data.Score = getScore(data)
	r.MergedData = append(r.MergedData, data)
}

func (r *Result) GetSortedData() []Data {
	r.sortData()
	return r.MergedData
}

func (r *Result) sortData() {
	sort.Slice(r.MergedData, func(i, j int) bool {
		return r.MergedData[i].Score > r.MergedData[j].Score
	})
}

func buildMetadata(data Data) map[string]string {
	metadata := make(map[string]string)
	for _, field := range config.Conf.Result.Score.MetadataFields {
		switch field {
		case "engine":
			metadata[field] = data.Engine
		case "title":
			metadata[field] = data.Title
		case "content":
			metadata[field] = data.Content
		case "$QUERY":
			metadata[field] = data.Query

		}
	}
	return metadata
}
