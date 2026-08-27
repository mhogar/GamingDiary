package cmds

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
)

type SeriesSelect struct {
	Index *int
}

func NewSeriesSelect(flags *flag.FlagSet) SeriesSelect {
	return SeriesSelect{
		Index: flags.Int("x", 0, "series index"),
	}
}

func (s SeriesSelect) Select(path string) (string, error) {
	data, err := json.UnmarshalFile[BaseData](filepath.Join(PATH, "index.json"))
	if err != nil {
		return "", errors.Chain(err, "error reading data file")
	}

	if *s.Index <= 0 || *s.Index > len(data.URLs) {
		for _, s := range data.URLs {
			fmt.Println(s)
		}
		return "", errors.Format("invalid index \"%d\"", *s.Index)
	}

	return data.URLs[*s.Index-1], nil
}
