package data

import "regexp"

var ENTRY_REGEX = regexp.MustCompile(`entry([0-9]+)\.json`)

const (
	STATIC_DIR        = "series"
	ENTRY_DATE_FORMAT = "Jan 02, 2006"
)
