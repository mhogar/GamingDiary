package data

import "regexp"

var ENTRY_REGEX = regexp.MustCompile(`entry([0-9]+)\.json`)

const (
	PUBLIC_PATH       = "public"
	STATIC_PATH       = "data/static"
	LOGS_PATH         = "data/logs"
	YOUTUBE_DATA_PATH = "data/youtube"
	TEMPLATE_PATH     = "data/templates"
)
