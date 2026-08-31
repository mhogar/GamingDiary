package data

type Root struct {
	Series []string `json:"series"`
}

type Series struct {
	//Index       string            `json:"index"`
	Title       string            `json:"title"`
	Dates       string            `json:"dates"`
	Description string            `json:"description"`
	Raw         string            `json:"raw"`
	Entries     string            `json:"entries"`
	Theme       string            `json:"theme"`
	Background  string            `json:"background"`
	Thumbnail   string            `json:"thumbnail"`
	Stylesheets []string          `json:"stylesheets"`
	Groups      map[string]string `json:"groups"`
}

type Entries struct {
	VideoCount    int     `json:"video_count"`
	TotalDuration float32 `json:"total_duration"`
	Entries       []Entry `json:"entries"`
}

type Entry struct {
	//Index       string  `json:"index"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float32 `json:"duration"`
	Thumbnail   string  `json:"thumbnail"`
	Video       string  `json:"video"`
	YouTube     string  `json:"youtube"`
}

type Parser interface {
	Parse(raw string) (Entry, error)
}
