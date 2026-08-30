package types

type Entry struct {
	Index       string  `json:"index"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float32 `json:"duration"`
	Thumbnail   string  `json:"thumbnail"`
	Video       string  `json:"video"`
	YouTube     string  `json:"youtube"`
}
