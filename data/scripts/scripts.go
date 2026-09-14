package scripts

import (
	"gamingdiary/data"

	"github.com/binarysoupdev/go-extensions/errors"
)

var scripts = []Scripts{
	YoutubeURLs{},
}

type Scripts interface {
	GetName() string
	Run(index string, entry *data.Entry) error
}

func Select(name string) (Scripts, error) {
	for _, s := range scripts {
		if s.GetName() == name {
			return s, nil
		}
	}
	return nil, errors.Format("invalid script \"%s\"", name)
}
