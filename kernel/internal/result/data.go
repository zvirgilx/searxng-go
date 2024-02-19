package result

import "regexp"

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

	// metadata stores data about fields and values in data.
	// This field is not shown externally and is only used internally.
	metadata map[string]string

	// score are scored by enabled Scorer on the data.
	score int
}

// unstructured converts the Data to a map.
func (d *Data) unstructured() *Data {
	metadata := make(map[string]string)
	for _, field := range conf.Score.MetadataFields {
		switch field {
		case "engine":
			metadata[field] = d.Engine
		case "title":
			metadata[field] = d.Title
		case "content":
			metadata[field] = d.Content
		case "$QUERY":
			metadata[field] = d.Query

		}
	}
	d.metadata = metadata
	return d
}

// match the variable in conditions values.
var variableMatcher = regexp.MustCompile(`^\$[A-Z_][A-Z0-9_]*$`)

// replace variable with real value.
func replaceVariable(origin []string, metadata map[string]string) {
	for i := range origin {
		// the variable value will replace by real value.
		// e.g. $QUERY -> query(query from search).
		if variableMatcher.MatchString(origin[i]) {
			origin[i] = metadata[origin[i]]
		}
	}
}

// doScore get a scorer and score the data.
func (d *Data) doScore() *Data {
	d.score = getScorer()(d.metadata)
	return d
}
