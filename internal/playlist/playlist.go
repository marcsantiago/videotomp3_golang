package playlist

type PlayList struct {
	Extractor    string `json:"extractor"`
	Type         string `json:"_type"`
	Title        string `json:"title"`
	ExtractorKey string `json:"extractor_key"`
	WebpageURL   string `json:"webpage_url"`
	Entries      []struct {
		URL   string `json:"url"`
		Type  string `json:"_type"`
		IeKey string `json:"ie_key"`
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"entries"`
	ID                 string `json:"id"`
	WebpageURLBasename string `json:"webpage_url_basename"`
}
