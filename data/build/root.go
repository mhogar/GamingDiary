package build

type Root struct {
	Background string                 `json:"background"`
	Logo       string                 `json:"logo"`
	Series     map[string]SeriesCache `json:"series"`
}
