package cmds

import (
	"fmt"
	"local/data"
	"local/templates"
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
}

type SeriesHeaderData struct {
	Path          string
	Title         string
	Dates         string
	Description   string
	VideoCount    int
	TotalDuration float32
	Thumbnail     string
	Theme         string
}

//=======================================

type RenderCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewRenderCommand() *RenderCommand {
	return &RenderCommand{
		CommandBase: command.NewCommandBase("render", "render the template"),
	}
}

func (cmd *RenderCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd RenderCommand) Run(args []string) error {
	dest := cmd.Flags.String("dest", "public", "the destination path")
	cmd.ParseFlags(args)

	root, err := json.UnmarshalFile[data.Root]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	homePage := HomePageData{
		Series: make([]SeriesHeaderData, len(root.Series)),
	}

	for i, name := range root.Series {
		style.Bold.Println(name)
		dataPath := filepath.Join("data", name)

		series, err := json.UnmarshalFile[data.Series](filepath.Join(dataPath, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading data file")
		}

		entries, err := json.UnmarshalFile[data.Entries](filepath.Join(dataPath, series.Entries))
		if err != nil {
			style.Error.Printf("x \"%s\" invalid or missing\n", series.Entries)
		}

		homePage.Series[i] = SeriesHeaderData{
			Path:          name,
			Title:         series.Title,
			Dates:         series.Dates,
			Description:   series.Description,
			VideoCount:    entries.VideoCount,
			TotalDuration: entries.TotalDuration,
			Thumbnail:     series.Thumbnail,
			Theme:         series.Theme,
		}
		homePage.VideoCount += entries.VideoCount
		homePage.TotalDuration += entries.TotalDuration

		resourcePath, _ := filepath.Rel(dataPath, "data")

		if err := cmd.renderSeries(series, entries.Entries, resourcePath, filepath.Join(*dest, name)); err != nil {
			return errors.ChainFormat(err, "error rendering series \"%s\"", name)
		}
	}

	style.Bold.Println("root")
	return cmd.renderHomePage(*dest, homePage)
}

func (cmd RenderCommand) renderHomePage(dest string, data HomePageData) error {
	page := templates.HomePage{
		VideoCount:    data.VideoCount,
		TotalDuration: cmd.formatDurationHMS(data.TotalDuration),
		Series:        make([]templates.SeriesHeader, len(data.Series)),
	}

	for i, header := range data.Series {
		page.Series[i] = templates.SeriesHeader{
			Title:         header.Title,
			Dates:         header.Dates,
			Description:   header.Description,
			VideoCount:    header.VideoCount,
			TotalDuration: cmd.formatDurationTimestamp(header.TotalDuration),
			Thumbnail:     filepath.Join(header.Path, header.Thumbnail),
			Link:          filepath.Join(header.Path, "index.html"),
			Theme:         header.Theme,
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

func (cmd *RenderCommand) renderSeries(series data.Series, entires []data.Entry, resourcePath, public string) error {
	page := templates.SeriesPage{
		Title:        series.Title,
		Dates:        series.Dates,
		Background:   series.Background,
		Theme:        series.Theme,
		Stylesheets:  series.Stylesheets,
		Entries:      make([]templates.Entry, len(entires)),
		ResourcePath: resourcePath,
	}

	for i, entry := range entires {
		page.Entries[i] = templates.Entry{
			Title:            entry.Title,
			Description:      entry.Description,
			Duration:         cmd.formatDurationTimestamp(entry.Duration),
			Thumbnail:        entry.Thumbnail,
			DefaultThumbnail: series.Thumbnail,
			Video:            entry.Video,
			Youtube:          entry.Youtube,
			Classes:          entry.Groups,
		}
	}

	out := filepath.Join(public, "index.html")

	err := os.MkdirAll(public, 0755)
	if err != nil {
		return errors.Chain(err, "error creating series path")
	}

	err = templates.RenderSeriesPage(out, page)
	if err != nil {
		return errors.Chain(err, "error rendering series page")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (RenderCommand) formatDurationTimestamp(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}

func (RenderCommand) formatDurationHMS(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%dh %dm %ds", d/(60*60), (d/60)%60, d%60)
}
