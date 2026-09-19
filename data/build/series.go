package build

import "time"

type Series struct {
	Index       int      `json:"index"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Theme       string   `json:"theme"`
	Background  string   `json:"background"`
	Thumbnail   string   `json:"thumbnail"`
	Stylesheets []string `json:"stylesheets"`
	SubSeries   []string `json:"sub_series"`
}

type SeriesCache struct {
	Index       int         `json:"index"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Thumbnail   string      `json:"thumbnail"`
	Theme       string      `json:"theme"`
	SubSeries   []string    `json:"sub_series"`
	Stats       SeriesStats `json:"stats"`
}

type SeriesStats struct {
	VideoCount    int       `json:"video_count"`
	TotalDuration float32   `json:"total_duration"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}
