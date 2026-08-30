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
}

//====================================================

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
	all := cmd.Flags.Bool("all", false, "render all templates")
	s := NewSeriesSelect(cmd.Flags)
	cmd.ParseFlags(args)

	if *all || *s.Index == 0 {
		if err := cmd.renderHomePage(*public); err != nil {
			return err
		}
	}

	if *all {
		return cmd.renderAllSeries(*public)
	} else if *s.Index != 0 {
		return cmd.renderSingleSeries(*public, s)
	}

	return nil
}

func (cmd RenderCommand) renderHomePage(public string) error {
	data, err := json.UnmarshalFile[BaseData]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	page := templates.HomePage{
		Series: make([]templates.SeriesHeader, len(data.Series)),
	}

	for i, series := range data.Series {
		dataPath := filepath.Join("data", series)

		s, err := json.UnmarshalFile[SeriesData](filepath.Join(dataPath, "index.json"))
		if err != nil {
			return errors.Chain(err, "error reading series data file")
		}

		page.Series[i] = templates.SeriesHeader{
			Title:       fmt.Sprintf("(%d) %s", i+1, s.Title),
			Dates:       s.Dates,
			Description: s.Description,
			//VideoCount:  len(durations),
			//TotalDuration: cmd.formatDuration(total),
			Thumbnail: filepath.Join(series, s.Thumbnail),
			Link:      filepath.Join(series, "index.html"),
			Theme:     s.Theme,
		}
	}

	out := filepath.Join(public, "index.html")

	err = templates.RenderHomePage(out, page)
	if err != nil {
		return errors.Chain(err, "error rendering home page")
	}

	style.Create.Printf("+ %s\n", out)
	return nil
}

func (cmd RenderCommand) renderAllSeries(public string) error {
	data, err := json.UnmarshalFile[BaseData]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	for _, series := range data.Series {
		if err := cmd.renderSeries(public, series); err != nil {
			return err
		}
	}

	return nil
}

func (cmd RenderCommand) renderSingleSeries(public string, s SeriesSelect) error {
	series, err := s.Select()
	if err != nil {
		return err
	}

	style.BoldInfo.Println(series)
	return cmd.renderSeries(public, series)
}

func (cmd RenderCommand) renderSeries(public, series string) error {
	dataPath := filepath.Join("data", series)

	data, err := json.UnmarshalFile[SeriesData](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

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

func (RenderCommand) formatDuration(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}
