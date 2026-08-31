package cmds

import (
	"fmt"
	"local/data"
	"local/templates"
	"math"
	"os"
	"path/filepath"
	"regexp"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type HomePageData struct {
	Series map[string]data.Series
	Keys   []string

	VideoCount    int
	TotalDuration float32
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
	public := cmd.Flags.String("public", "public", "the public path")
	cmd.ParseFlags(args)

	root, err := json.UnmarshalFile[data.Root]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	homePage := HomePageData{
		Series: map[string]data.Series{},
		Keys:   make([]string, len(root.Series)),
	}

	for i, name := range root.Series {
		dataPath := filepath.Join("data", name)

		series, err := json.UnmarshalFile[data.Series](filepath.Join(dataPath, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading data file")
		}

		homePage.Series[name] = series
		homePage.Keys[i] = name

		entires, err := json.UnmarshalFile[data.Entries](filepath.Join(dataPath, series.Entries))
		if err != nil {
			return errors.Chain(err, "error reading entries file")
		}

		homePage.VideoCount++
		homePage.TotalDuration += entires.TotalDuration

		resourcePath, _ := filepath.Rel(dataPath, "data")

		if err := cmd.renderSeries(series, entires.Entries, resourcePath, filepath.Join(*public, name)); err != nil {
			return errors.ChainFormat(err, "error rendering series \"%s\"", name)
		}
	}

	return cmd.renderHomePage(*public, homePage)
}

func (cmd *RenderCommand) renderSeries(series data.Series, entires []data.Entry, resourcePath, public string) error {
	groups := make(map[string]*regexp.Regexp)
	for key, val := range series.Groups {
		groups[key] = regexp.MustCompile(val)
	}

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
		classes := []string{}
		for group, regex := range groups {
			if regex.MatchString(entry.Index) {
				classes = append(classes, group)
			}
		}

		page.Entries[i] = templates.Entry{
			Title:       entry.Title,
			Description: entry.Description,
			Duration:    cmd.formatDuration(entry.Duration),
			Thumbnail:   entry.Thumbnail,
			Video:       entry.Video,
			Classes:     classes,
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

func (cmd RenderCommand) renderHomePage(public string, data HomePageData) error {
	page := templates.HomePage{
		VideoCount:    data.VideoCount,
		TotalDuration: cmd.formatDuration(data.TotalDuration),
		Series:        make([]templates.SeriesHeader, len(data.Series)),
	}

	for i, name := range data.Keys {
		series := data.Series[name]

		page.Series[i] = templates.SeriesHeader{
			Title:         fmt.Sprintf("(%s) %s", series.Index, series.Title),
			Dates:         series.Dates,
			Description:   series.Description,
			VideoCount:    data.VideoCount,
			TotalDuration: cmd.formatDuration(data.TotalDuration),
			Thumbnail:     filepath.Join(name, series.Thumbnail),
			Link:          filepath.Join(name, "index.html"),
			Theme:         series.Theme,
		}
	}

	out := filepath.Join(public, "index.html")

	err := templates.RenderHomePage(out, page)
	if err != nil {
		return errors.Chain(err, "error rendering home page")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (RenderCommand) formatDuration(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}
