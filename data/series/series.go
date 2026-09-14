package series

import (
	"gamingdiary/data"
	"gamingdiary/tools/youtube"

	"github.com/binarysoupdev/go-extensions/errors"
)

var series = []Series{
	SunshineChapters{}, SunshineShorts{},
	TTYDSeries{}, TTYDBattles{}, TTYDShorts{},
	LuigiMansionSeries{}, LuigiMansionShorts{},
	ShakeItSeries{},
	HeartgoldSeries{},
}

type Series interface {
	GetName() string
	BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error
}

func Select(name string) (Series, error) {
	for _, s := range series {
		if s.GetName() == name {
			return s, nil
		}
	}
	return nil, errors.Format("invalid series \"%s\"", name)
}
