package series

import (
	"app/data"
	"app/tools/youtube"

	"github.com/binarysoupdev/go-extensions/errors"
)

type Series interface {
	GetName() string
	BuildNewEntry(index int, entry *data.Entry) error
	BuildEntryFromVideo(index int, video string, entry *data.Entry) error
	BuildEntryFromYoutube(index int, video *youtube.Video, entry *data.Entry) error
}

type seriesBase struct{}

func (seriesBase) BuildNewEntry(_ int, _ *data.Entry) error {
	return nil
}

func (seriesBase) BuildEntryFromVideo(_ int, _ string, _ *data.Entry) error {
	return nil
}

func (seriesBase) BuildEntryFromYoutube(_ int, _ *youtube.Video, _ *data.Entry) error {
	return nil
}

//=================================================

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

func Select(name string) (Series, error) {
	for _, s := range series {
		if s.GetName() == name {
			return s, nil
		}
	}
	return nil, errors.Format("invalid series \"%s\"", name)
}
