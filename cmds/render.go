package cmds

import (
	"fmt"
	"local/cmds/types"
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

type BaseData struct {
	Series []string `json:"series"`
}

type SeriesData struct {
	Title       string            `json:"title"`
	Dates       string            `json:"dates"`
	Description string            `json:"description"`
	RawFiles    string            `json:"raw_files"`
	Entries     string            `json:"entries"`
	Theme       string            `json:"theme"`
	Background  string            `json:"background"`
	Thumbnail   string            `json:"thumbnail"`
	Stylesheets []string          `json:"stylesheets"`
	Groups      map[string]string `json:"groups"`

	VideoCount    int     `json:"video_count"`
	TotalDuration float32 `json:"total_duration"`
}

//====================================================

type RenderCommand struct {
	command.CommandBase
	command.FlagCommand

	videoCount    int
	totalDuration float32
}

func NewRenderCommand() *RenderCommand {
	return &RenderCommand{
		CommandBase: command.NewCommandBase("render", "render the template"),
	}
}

func (cmd *RenderCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)

	cmd.videoCount = 0
	cmd.totalDuration = 0

	return nil
}

func (cmd RenderCommand) Run(args []string) error {
	public := cmd.Flags.String("public", "public", "the public path")
	cmd.ParseFlags(args)

	data, err := json.UnmarshalFile[BaseData]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	for _, series := range data.Series {
		if err := cmd.renderSeries(*public, series); err != nil {
			return err
		}
	}

	return cmd.renderHomePage(*public, data)
}

func (cmd *RenderCommand) renderSeries(public, series string) error {
	dataPath := filepath.Join("data", series)

	data, err := json.UnmarshalFile[SeriesData](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	cmd.videoCount += data.VideoCount
	cmd.totalDuration += data.TotalDuration

	groups := make(map[string]*regexp.Regexp)
	for key, val := range data.Groups {
		groups[key] = regexp.MustCompile(val)
	}

	files, err := filepath.Glob(filepath.Join(dataPath, data.Entries))
	if err != nil {
		return errors.Chain(err, "error finding entries")
	}

	page := templates.SeriesPage{
		Title:       data.Title,
		Dates:       data.Dates,
		Background:  data.Background,
		Theme:       data.Theme,
		Stylesheets: data.Stylesheets,
		Entries:     make([]templates.Entry, len(files)),
	}

	for i, file := range files {
		entry, err := json.UnmarshalFile[types.Entry](file)
		if err != nil {
			return errors.Chain(err, "error reading entry file")
		}

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

	seriesPath := filepath.Join(public, series)
	out := filepath.Join(seriesPath, "index.html")

	err = os.MkdirAll(seriesPath, 0755)
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

func (cmd RenderCommand) renderHomePage(public string, data BaseData) error {
	page := templates.HomePage{
		VideoCount:    cmd.videoCount,
		TotalDuration: cmd.formatDuration(cmd.totalDuration),
		Series:        make([]templates.SeriesHeader, len(data.Series)),
	}

	for i, series := range data.Series {
		dataPath := filepath.Join("data", series)

		data, err := json.UnmarshalFile[SeriesData](filepath.Join(dataPath, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading series data file")
		}

		page.Series[i] = templates.SeriesHeader{
			Title:         fmt.Sprintf("(%d) %s", i+1, data.Title),
			Dates:         data.Dates,
			Description:   data.Description,
			VideoCount:    data.VideoCount,
			TotalDuration: cmd.formatDuration(data.TotalDuration),
			Thumbnail:     filepath.Join(series, data.Thumbnail),
			Link:          filepath.Join(series, "index.html"),
			Theme:         data.Theme,
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
