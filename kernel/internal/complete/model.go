package complete

type Result struct {
	Type string `json:"type"` // Type determines whether it is a text-only type or a media type.
	Text string `json:"text"` // Text is the associative suggestion text for query.
	Info string `json:"info"` // Info contains the media information required by the media type.
}
