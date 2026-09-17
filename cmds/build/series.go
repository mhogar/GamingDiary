package build_cmd

import (
	"app/data"
	"app/data/templates"
	"fmt"
	"os"
	"path/filepath"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
)

func (cmd BuildCommand) buildSeries(dest, name string) (data.SeriesCache, error) {
	seriesPath := filepath.Join(data.STATIC_PATH, name)

	series, err := json.UnmarshalFile[data.Series](filepath.Join(seriesPath, "series.json"))
	if err != nil {
		cmd.logError(err, "invalid series file")
		return data.SeriesCache{}, errors.New("invalid series file")
	}

	cache := data.SeriesCache{
		Index:       series.Index,
		Title:       series.Title,
		Description: series.Description,
		Thumbnail:   series.Thumbnail,
		Theme:       series.Theme,
		SubSeries:   make([]string, 0, len(series.SubSeries)),
	}

	for i, subSeries := range series.SubSeries {
		stats, err := cmd.buildSubSeries(filepath.Join(dest, name, subSeries), name, subSeries, series)
		if err != nil {
			cmd.printError(err.Error())
			continue
		}

		cache.SubSeries = append(cache.SubSeries, subSeries)
		cache.Stats.VideoCount += stats.VideoCount
		cache.Stats.TotalDuration += stats.TotalDuration

		if i == 0 {
			cache.Stats.StartDate = stats.StartDate
			cache.Stats.EndDate = stats.EndDate
		}
	}

	cmd.printSeriesHeader(name)
	fs := fileStats{}

	files := make([]string, 0, 2+len(series.Stylesheets))
	files = append(files, series.Background, series.Thumbnail)
	files = append(files, series.Stylesheets...)

	cmd.copyFiles(filepath.Join(dest, name), filepath.Join(data.PUBLIC_PATH, name), &fs, files)
	fs.Print()

	return cache, nil
}

func (cmd *BuildCommand) buildSubSeries(dest, seriesName, subSeries string, series data.Series) (data.SeriesStats, error) {
	name := filepath.Join(seriesName, subSeries)
	cmd.logBuild(name)
	cmd.printSeriesHeader(name)

	entries, err := filepath.Glob(filepath.Join(data.STATIC_PATH, name, "entry*.json"))
	if err != nil {
		cmd.logError(err, "error reading sub-series directory")
		return data.SeriesStats{}, errors.New("error reading directory")
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		cmd.logError(err, "error creating sub-series directory")
		return data.SeriesStats{}, errors.New("error creating out directory")
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
	fs := fileStats{}

	for i, file := range entries {
		entry, err := json.UnmarshalFile[data.Entry](file)
		if err != nil {
			cmd.logError(err, fmt.Sprintf("error reading entry file \"%s\"", file))
			cmd.printError("invalid " + filepath.Base(file))
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

		if cmd.local {
			files := []string{entry.Thumbnail, entry.Video}
			cmd.copyFiles(dest, filepath.Join(data.PUBLIC_PATH, name), &fs, files)
		}
	}

	page.Dates = cmd.formatDateRange(stats.StartDate, stats.EndDate)
	page.TotalDuration = cmd.formatDurationTimestamp(stats.TotalDuration)

	out := filepath.Join(dest, "index.html")

	if err := templates.RenderSeriesPage(out, page); err != nil {
		cmd.logError(err, "error rendering series page")
		return stats, errors.Chain(err, "error rendering page")
	}

	cmd.logCreate(out)
	cmd.printCreate(out)

	fs.Created++
	fs.Print()

	return stats, nil
}
