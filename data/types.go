package data

type Root struct {
	Series []string `json:"series"`
}

type Series struct {
	Title           string            `json:"title"`
	Dates           string            `json:"dates"`
	Description     string            `json:"description"`
	Entries         string            `json:"entries"`
	Theme           string            `json:"theme"`
	Background      string            `json:"background"`
	Thumbnail       string            `json:"thumbnail"`
	Stylesheets     []string          `json:"stylesheets"`
	Groups          map[string]string `json:"groups"`
	YoutubePlaylist string            `json:"youtube_playlist"`
}

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
	Youtube     string   `json:"youtube"`
	Groups      []string `json:"groups"`
}
