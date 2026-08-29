package cmds

import (
	"flag"
	"fmt"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
)

type SeriesSelect struct {
	Index *int
}

func NewSeriesSelect(flags *flag.FlagSet) SeriesSelect {
	return SeriesSelect{
		Index: flags.Int("x", -1, "series index"),
	}
}

func (s SeriesSelect) Select() (string, error) {
	data, err := json.UnmarshalFile[BaseData]("data/index.json")
	if err != nil {
		return "", errors.Chain(err, "error reading data file")
	}

	if *s.Index <= 0 || *s.Index > len(data.Series) {
		for i, s := range data.Series {
			fmt.Printf("[%d] %s\n", i+1, s)
		}
		return "", errors.Format("invalid index \"%d\"", *s.Index)
	}

	return data.Series[*s.Index-1], nil
}
