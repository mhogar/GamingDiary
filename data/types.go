package data

type Root struct {
	Series []string `json:"series"`
}

type Series struct {
	Title       string      `json:"title"`
	Dates       string      `json:"dates"`
	Description string      `json:"description"`
	Theme       string      `json:"theme"`
	Background  string      `json:"background"`
	Thumbnail   string      `json:"thumbnail"`
	Stylesheets []string    `json:"stylesheets"`
	SubSeries   []SubSeries `json:"sub_series"`
}

type SubSeries struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// Deprecated
type Entries struct {
	VideoCount    int     `json:"video_count"`
	TotalDuration float32 `json:"total_duration"`
	Entries       []Entry `json:"entries"`
}

type Entry struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Duration    float32  `json:"duration"`
	Thumbnail   string   `json:"thumbnail"`
	Video       string   `json:"video"`
	Youtube     string   `json:"youtube,omitempty"` // Deprecated
	YoutubeId   string   `json:"youtube_id"`
	Groups      []string `json:"groups"`
}
