package build

type Root struct {
	Icon       string `json:"icon"`
	Background string `json:"background"`
}

type SeriesMap map[string]SeriesCache
