package cmds

import (
	"flag"
	"fmt"
	"os"

	"github.com/binarysoupdev/go-extensions/errors"
)

type SeriesSelect struct {
	index *int
}

func NewSeriesSelect(flags *flag.FlagSet) SeriesSelect {
	return SeriesSelect{
		index: flags.Int("x", 0, "series index"),
	}
}

func (s SeriesSelect) Select(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", errors.Chain(err, "error reading directory")
	}
	series := make([]string, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			series = append(series, file.Name())
		}
	}

	if *s.index <= 0 || *s.index > len(series) {
		for _, s := range series {
			fmt.Println(s)
		}
		return "", errors.Format("invalid index \"%d\"", *s.index)
	}

	return series[*s.index-1], nil
}
