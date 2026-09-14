package series

import (
	"app/data"
	"app/tools/youtube"

	"github.com/binarysoupdev/go-extensions/errors"
)

var series = []Series{
	SunshineChapters{}, SunshineShorts{},
	TTYDSeries{}, TTYDBattles{}, TTYDShorts{},
	LuigiMansionSeries{}, LuigiMansionShorts{},
	ShakeItSeries{},
	WinniePoohSeries{},
	HeartgoldSeries{},
	OrigamiKingSeries{},
	LuigiMansion3Series{},
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
