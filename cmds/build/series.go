package build_cmd

import (
	"app/data"
	"app/data/templates"
	"fmt"
	"os"
	"path/filepath"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func (cmd BuildCommand) selectSeries(name string, root data.Root) ([]string, error) {
	switch name {
	case "root":
		return []string{}, nil
	case "all":
		s := make([]string, 0, len(root.Series))
		for name := range root.Series {
			s = append(s, name)
		}
		return s, nil
	default:
		stat, err := os.Stat(filepath.Join(data.STATIC_PATH, name))
		if err != nil || !stat.IsDir() {
			return nil, errors.Format("invalid series \"%s\"", name)
		}
		return []string{name}, nil
	}
}

func (cmd BuildCommand) buildSeries(dest, name string) (data.SeriesCache, error) {
	style.Bold.Println(name)
	seriesPath := filepath.Join(data.STATIC_PATH, name)

	s, err := json.UnmarshalFile[data.Series](filepath.Join(seriesPath, "series.json"))
	if err != nil {
		cmd.logError(err, "invalid series file")
		return data.SeriesCache{}, errors.New("invalid series file")
	}

	series := data.SeriesCache{
		Index:       s.Index,
		Title:       s.Title,
		Description: s.Description,
		Thumbnail:   s.Thumbnail,
		Theme:       s.Theme,
		SubSeries:   make([]string, 0, len(s.SubSeries)),
	}

	for i, subSeries := range s.SubSeries {
		stats, err := cmd.buildSubSeries(name, subSeries, s, filepath.Join(dest, name, subSeries))
		if err != nil {
			style.Error.Printf("[x] (%s) %s\n", subSeries, err)
			continue
		}

		series.SubSeries = append(series.SubSeries, subSeries)
		series.Stats.VideoCount += stats.VideoCount
		series.Stats.TotalDuration += stats.TotalDuration

		if i == 0 {
			series.Stats.StartDate = stats.StartDate
			series.Stats.EndDate = stats.EndDate
		}
	}

	// files := make([]string, 0, 2+len(s.Stylesheets))
	// files = append(files, s.Background, s.Thumbnail)
	// files = append(files, s.Stylesheets...)
	// cmd.copyFiles(filepath.Join(dest, name), filepath.Join(data.STATIC_PATH, name), files)

	return series, nil
}

func (cmd *BuildCommand) buildSubSeries(seriesName, subSeries string, series data.Series, dest string) (data.SeriesStats, error) {
	out := filepath.Join(dest, "index.html")

	entries, err := filepath.Glob(filepath.Join(data.STATIC_PATH, seriesName, subSeries, "entry*.json"))
	if err != nil {
		return data.SeriesStats{}, errors.Chain(err, "error reading directory")
	}

	if len(entries) == 0 {
		return data.SeriesStats{}, errors.New("no entries found")
	}

	page := templates.SeriesPage{
		Title:       fmt.Sprintf("%s (%s)", series.Title, subSeries),
		Background:  series.Background,
		Theme:       series.Theme,
		Stylesheets: series.Stylesheets,
		Entries:     make([]templates.Entry, len(entries)),
	}

	stats := data.SeriesStats{
		VideoCount: len(entries),
	}

	var errs errors.Errors
	for i, file := range entries {
		entry, err := json.UnmarshalFile[data.Entry](file)
		if err != nil {
			errs.Add(err)
			continue
		}

		stats.TotalDuration += entry.Duration

		page.Entries[i] = templates.Entry{
			Title:            entry.Title,
			Description:      entry.Description,
			Duration:         cmd.formatDurationTimestamp(entry.Duration),
			Date:             entry.Date.Format(data.DATE_FORMAT),
			Thumbnail:        entry.Thumbnail,
			DefaultThumbnail: filepath.Join("..", series.Thumbnail),
			Video:            entry.Video,
			YoutubeThumbnail: entry.YoutubeThumbnail,
			YoutubeVideo:     entry.YoutubeVideo,
			Group:            entry.Group,
			Local:            cmd.local,
		}

		if i == 0 {
			stats.StartDate = entry.Date
		} else if i == len(entries)-1 {
			stats.EndDate = entry.Date
		}
	}

	if len(errs) > 0 {
		return stats, errors.Chain(errs.Collapse("\n  "), "error building entry")
	}

	page.Dates = cmd.formatDateRange(stats.StartDate, stats.EndDate)
	page.TotalDuration = cmd.formatDurationTimestamp(stats.TotalDuration)

	if err := os.MkdirAll(dest, 0755); err != nil {
		return stats, errors.Chain(err, "error creating series path")
	}

	if err := templates.RenderSeriesPage(out, page); err != nil {
		return stats, errors.Chain(err, "error rendering series page")
	}

	style.Create.Printf("+ %s\n", out)
	return stats, nil
}
