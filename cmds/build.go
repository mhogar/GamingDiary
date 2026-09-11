package cmds

import (
	"fmt"
	"gamingdiary/data"
	"gamingdiary/templates"
	"gamingdiary/util"
	"math"
	"os"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type HomePageData struct {
	Series        []SeriesHeaderData
	VideoCount    int
	TotalDuration float32
	StartDate     string
}

type SeriesHeaderData struct {
	Path          string
	Title         string
	StartDate     string
	EndDate       string
	Description   string
	VideoCount    int
	TotalDuration float32
	Thumbnail     string
	Theme         string
	SubSeries     []data.SubSeries
}

type SubSeriesStats struct {
	VideoCount    int
	TotalDuration float32
	StartDate     string
	EndDate       string
}

//=======================================

type BuildCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewBuildCommand() *BuildCommand {
	return &BuildCommand{
		CommandBase: command.NewCommandBase("build", "Build the app from the templates"),
	}
}

func (cmd *BuildCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd BuildCommand) Run(args []string) error {
	//series := cmd.Flags.String("series", "", "name of the series")
	dest := cmd.Flags.String("dest", "series", "the destination path")
	cmd.ParseFlags(args)

	// if *series == "" {
	// 	return errors.New("\"series\" cannot be empty")
	// }
	// style.BoldInfo.Println(*series)

	root, err := json.UnmarshalFile[data.Root]("series/root.json")
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	homePage := HomePageData{
		Series: make([]SeriesHeaderData, len(root.Series)),
	}

	for i, name := range root.Series {
		style.Bold.Println(name)
		seriesPath := filepath.Join("series", name)

		series, err := json.UnmarshalFile[data.Series](filepath.Join(seriesPath, "series.json"))
		if err != nil {
			return errors.Chain(err, "error reading series file")
		}

		header := SeriesHeaderData{
			Path:        name,
			Title:       series.Title,
			Description: series.Description,
			Thumbnail:   series.Thumbnail,
			Theme:       series.Theme,
			SubSeries:   series.SubSeries,
		}

		for i, subSeries := range series.SubSeries {
			stats, err := cmd.renderSubSeries(series, subSeries, filepath.Join(seriesPath, subSeries.Path), filepath.Join(*dest, name, subSeries.Path))
			if err != nil {
				return errors.ChainFormat(err, "error rendering series \"%s\"", name)
			}

			header.VideoCount += stats.VideoCount
			header.TotalDuration += stats.TotalDuration

			if i == 0 {
				header.StartDate = stats.StartDate
				header.EndDate = stats.EndDate
			}
		}

		homePage.Series[i] = header
		homePage.VideoCount += header.VideoCount
		homePage.TotalDuration += header.TotalDuration
	}
	homePage.StartDate = homePage.Series[0].StartDate

	style.Bold.Println("root")
	return cmd.renderHomePage(*dest, homePage)
}

func (cmd BuildCommand) renderHomePage(dest string, data HomePageData) error {
	page := templates.HomePage{
		VideoCount:    data.VideoCount,
		TotalDuration: cmd.formatDurationHMS(data.TotalDuration),
		StartDate:     data.StartDate,
		Series:        make([]templates.SeriesHeader, len(data.Series)),
	}

	for i, header := range data.Series {
		links := make([]templates.SubSeriesLink, len(header.SubSeries))
		for i, subSeries := range header.SubSeries {
			links[i] = templates.SubSeriesLink{
				Title:     util.Capitalize(subSeries.Title),
				Link:      filepath.Join(header.Path, subSeries.Path, "index.html"),
				Separator: " | ",
			}
		}
		links[len(links)-1].Separator = ""

		page.Series[i] = templates.SeriesHeader{
			Title:          fmt.Sprintf("(%d) %s", i+1, header.Title),
			Dates:          cmd.formatDateRange(header.StartDate, header.EndDate),
			Description:    header.Description,
			VideoCount:     header.VideoCount,
			TotalDuration:  cmd.formatDurationTimestamp(header.TotalDuration),
			Thumbnail:      filepath.Join(header.Path, header.Thumbnail),
			SubSeriesLinks: links,
			Theme:          header.Theme,
		}
	}

	out := filepath.Join(dest, "index.html")

	err := templates.RenderHomePage(out, page)
	if err != nil {
		return errors.Chain(err, "error rendering home page")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (cmd *BuildCommand) renderSubSeries(series data.Series, subSeries data.SubSeries, path, public string) (SubSeriesStats, error) {
	entries, err := filepath.Glob(filepath.Join(path, "entry*.json"))
	if err != nil {
		return SubSeriesStats{}, errors.Chain(err, "error reading series directory")
	}

	if len(entries) == 0 {
		return SubSeriesStats{}, errors.New("no entries found")
	}
	resPath, _ := filepath.Rel(path, "series")

	page := templates.SeriesPage{
		Title:        fmt.Sprintf("%s (%s)", series.Title, subSeries.Title),
		Background:   filepath.Join("..", series.Background),
		Theme:        series.Theme,
		Stylesheets:  series.Stylesheets,
		Entries:      make([]templates.Entry, len(entries)),
		ResourcePath: resPath,
	}

	stats := SubSeriesStats{
		VideoCount: len(entries),
	}

	for i, file := range entries {
		entry, err := json.UnmarshalFile[data.Entry](file)
		if err != nil {
			return stats, errors.Chain(err, "error reading entry")
		}

		stats.TotalDuration += entry.Duration

		page.Entries[i] = templates.Entry{
			Title:            entry.Title,
			Description:      entry.Description,
			Duration:         cmd.formatDurationTimestamp(entry.Duration),
			Date:             entry.Date,
			Thumbnail:        entry.Thumbnail,
			DefaultThumbnail: filepath.Join("..", series.Thumbnail),
			Video:            entry.Video,
			Youtube:          fmt.Sprintf("https://www.youtube.com/watch?v=%s", entry.YoutubeId),
			Classes:          entry.Groups,
		}
	}

	stats.StartDate = page.Entries[0].Date
	stats.EndDate = page.Entries[len(page.Entries)-1].Date

	page.Dates = cmd.formatDateRange(stats.StartDate, stats.EndDate)
	page.TotalDuration = cmd.formatDurationTimestamp(stats.TotalDuration)

	out := filepath.Join(public, "index.html")

	err = os.MkdirAll(public, 0755)
	if err != nil {
		return stats, errors.Chain(err, "error creating series path")
	}

	err = templates.RenderSeriesPage(out, page)
	if err != nil {
		return stats, errors.Chain(err, "error rendering series page")
	}

	style.Create.Printf("+ %s\n", out)
	return stats, nil
}

func (BuildCommand) formatDurationTimestamp(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDurationHMS(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%dh %dm %ds", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDateRange(startDate, endDate string) string {
	return fmt.Sprintf("%s - %s", startDate, endDate)
}
