package scripts

import (
	"app/data/build"

	"github.com/binarysoupdev/go-extensions/errors"
)

var scripts = []Scripts{}

type Scripts interface {
	GetName() string
	Run(index string, entry *build.Entry) error
}

func Select(name string) (Scripts, error) {
	for _, s := range scripts {
		if s.GetName() == name {
			return s, nil
		}
	}
	return nil, errors.Format("invalid script \"%s\"", name)
}
