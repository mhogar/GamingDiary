package data

import "time"

type Root struct {
	Background string                 `json:"background"`
	Logo       string                 `json:"logo"`
	Series     map[string]SeriesCache `json:"series"`
}

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

type Entry struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Date             time.Time `json:"date"`
	Duration         float32   `json:"duration"`
	Thumbnail        string    `json:"thumbnail"`
	Video            string    `json:"video"`
	YoutubeId        string    `json:"youtube_id"`
	YoutubeThumbnail string    `json:"youtube_thumbnail"`
	YoutubeVideo     string    `json:"youtube_video"`
	Groups           []string  `json:"groups,omitempty"` // Deprecated
	Group            string    `json:"group"`
}

type YoutubeMeta struct {
	Playlist string `json:"playlist"`
}
