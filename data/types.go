package data

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

// Deprecated
type Entries struct {
	VideoCount    int     `json:"video_count"`
	TotalDuration float32 `json:"total_duration"`
	Entries       []Entry `json:"entries"`
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
