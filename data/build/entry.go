package build

import "time"

type Entry struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Date             time.Time `json:"date"`
	Duration         float32   `json:"duration"`
	Thumbnail        string    `json:"thumbnail"`
	Video            string    `json:"video"`
	YoutubeId        string    `json:"youtube_id"`
	YoutubeThumbnail string    `json:"youtube_thumbnail,omitempty"`
	Group            string    `json:"group,omitempty"`
}
