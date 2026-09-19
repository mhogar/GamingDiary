package build_cmd

import (
	"app/data"
	"app/data/build"
	"app/data/templates"
	"app/util"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/binarysoupdev/go-extensions/errors"
)

func (cmd BuildCommand) buildRoot(dest, appName string, root build.Root) error {
	cmd.logBuild("root")
	cmd.printSeriesHeader("root")

	page := templates.RootPage{
		AppName:    util.Capitalize(appName),
		Logo:       root.Logo,
		Background: root.Background,
		VideoCount: 0,
		Series:     make([]templates.SeriesHeader, 0, len(root.Series)),
	}

	var duration float32
	var startDate time.Time
	var endDate time.Time

	for name, series := range root.Series {
		tmpl := templates.SeriesHeader{
			Index:          series.Index,
			Title:          series.Title,
			Dates:          cmd.formatDateRange(series.Stats.StartDate, series.Stats.EndDate),
			Description:    series.Description,
			VideoCount:     series.Stats.VideoCount,
			TotalDuration:  cmd.formatDurationTimestamp(series.Stats.TotalDuration),
			Thumbnail:      filepath.Join(name, series.Thumbnail),
			SubSeriesLinks: cmd.buildSubSeriesLinks(name, series),
			Theme:          series.Theme,
		}

		page.Series = append(page.Series, tmpl)
		page.VideoCount += series.Stats.VideoCount
		duration += series.Stats.TotalDuration

		if startDate.IsZero() || series.Stats.StartDate.Before(startDate) {
			startDate = series.Stats.StartDate
		}
		if endDate.IsZero() || series.Stats.EndDate.After(endDate) {
			endDate = series.Stats.EndDate
		}
	}

	page.TotalDuration = cmd.formatDurationHMS(duration)
	if len(page.Series) > 0 {
		page.Dates = fmt.Sprintf("%s - %s", startDate.Format(data.DATE_FORMAT), endDate.Format(data.DATE_FORMAT))
	}

	slices.SortFunc(page.Series, func(a, b templates.SeriesHeader) int {
		return a.Index - b.Index
	})

	out := filepath.Join(dest, "index.html")

	if err := templates.RenderRootPage(out, page); err != nil {
		cmd.logError(err, "error rendering root")
		return errors.New("error rendering root")
	}

	cmd.logCreate(out)
	cmd.printCreate(out)

	fs := fileStats{Created: 1}
	cmd.copyFiles(dest, data.PUBLIC_PATH, &fs, "style.css", "script.js", root.Background, root.Logo)
	fs.Print()

	return nil
}

func (cmd BuildCommand) buildSubSeriesLinks(name string, series build.SeriesCache) []templates.SubSeriesLink {
	if len(series.SubSeries) == 0 {
		return nil
	}
	links := make([]templates.SubSeriesLink, len(series.SubSeries))

	for i, subSeries := range series.SubSeries {
		links[i] = templates.SubSeriesLink{
			Title:     util.Capitalize(subSeries),
			Link:      filepath.Join(name, subSeries, "index.html"),
			Separator: " | ",
		}
	}

	links[len(links)-1].Separator = ""
	return links
}
