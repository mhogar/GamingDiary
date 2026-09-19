package data

import "regexp"

var ENTRY_REGEX = regexp.MustCompile(`entry([0-9]+)\.json`)
var URL_REGEX = regexp.MustCompile(`^https?://`)
var YOUTUBE_DURATION_REGEX = regexp.MustCompile(`([0-9]+)([^0-9])`)
var VIDEO_INDEX_REGEX = regexp.MustCompile(`v(.+)\.mp4`)

const (
	PUBLIC_PATH       = "public"
	STATIC_PATH       = "data/static"
	LOGS_PATH         = "data/logs"
	YOUTUBE_DATA_PATH = "data/youtube"
	TEMPLATE_PATH     = "data/templates"
	CONFIG_PATH       = "config.json"
	SERIES_CACHE_PATH = "data/cache/series.json"
)

const (
	DATE_FORMAT       = "Jan 02, 2006"
	YOUTUBE_META_FILE = "youtube.json"
	ENTRY_PATTERN     = "entry*.json"
)
