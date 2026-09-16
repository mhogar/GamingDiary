package cmds

import (
	"app/data"
	"app/data/templates"
	"app/util"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/file"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

const DATE_FORMAT = "Jan 02, 2006"

type RootData struct {
	Background string                `json:"background"`
	Logo       string                `json:"logo"`
	Series     map[string]SeriesData `json:"series"`
}

type SeriesData struct {
	Index       int         `json:"index"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Thumbnail   string      `json:"thumbnail"`
	Theme       string      `json:"theme"`
	SubSeries   []string    `json:"sub_series"`
	Stats       SeriesStats `json:"stats"`
}

type SeriesStats struct {
	VideoCount    int       `json:"video_count"`
	TotalDuration float32   `json:"total_duration"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

type FileStats struct {
	Copied   int
	UpToDate int
	NotFound int
	Invalid  int
}

//=======================================

type BuildCommand struct {
	command.CommandBase
	command.FlagCommand

	local  bool
	logger *log.Logger
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
	s := cmd.Flags.String("series", "", "name of the series")
	out := cmd.Flags.String("out", data.PUBLIC_PATH, "the destination path")
	cmd.Flags.BoolVar(&cmd.local, "local", false, "build using local thumbnails and videos")
	cmd.ParseFlags(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}
	rootPath := filepath.Join(data.STATIC_PATH, "root.json")

	root, err := json.UnmarshalFile[RootData](rootPath)
	if err != nil {
		return errors.Chain(err, "error reading root file")
	}

	series, err := cmd.selectSeries(*s, root)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(data.LOGS_PATH, fmt.Sprintf("build-%s.txt", time.Now().Format(time.DateTime))))
	if err != nil {
		return errors.Chain(err, "error creating log file")
	}
	defer f.Close()
	cmd.logger = log.New(f, "", log.Ltime)

	for _, s := range series {
		data, err := cmd.buildSeries(*out, s)
		if err == nil {
			root.Series[s] = data
		} else {
			style.Error.Printf("[x] %s\n", err)
		}
	}

	if err := cmd.buildRoot(*out, root); err != nil {
		return err
	}

	if err := json.MarshalFilePretty(root, rootPath, "    "); err != nil {
		return errors.Chain(err, "error saving root file")
	}
	return nil
}

func (cmd BuildCommand) selectSeries(name string, root RootData) ([]string, error) {
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

func (cmd BuildCommand) buildSeries(dest, name string) (SeriesData, error) {
	style.Bold.Println(name)
	seriesPath := filepath.Join(data.STATIC_PATH, name)

	s, err := json.UnmarshalFile[data.Series](filepath.Join(seriesPath, "series.json"))
	if err != nil {
		cmd.logError(err, "invalid series file")
		return SeriesData{}, errors.New("invalid series file")
	}

	series := SeriesData{
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

func (cmd *BuildCommand) buildSubSeries(seriesName, subSeries string, series data.Series, dest string) (SeriesStats, error) {
	out := filepath.Join(dest, "index.html")

	entries, err := filepath.Glob(filepath.Join(data.STATIC_PATH, seriesName, subSeries, "entry*.json"))
	if err != nil {
		return SeriesStats{}, errors.Chain(err, "error reading directory")
	}

	if len(entries) == 0 {
		return SeriesStats{}, errors.New("no entries found")
	}

	page := templates.SeriesPage{
		Title:       fmt.Sprintf("%s (%s)", series.Title, subSeries),
		Background:  series.Background,
		Theme:       series.Theme,
		Stylesheets: series.Stylesheets,
		Entries:     make([]templates.Entry, len(entries)),
	}

	stats := SeriesStats{
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
			Date:             entry.Date.Format(DATE_FORMAT),
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

	cmd.printCreate(out)
	return stats, nil
}

func (cmd BuildCommand) buildRoot(dest string, root RootData) error {
	cmd.logBuild("ROOT")
	style.New(style.BOLD, style.UNDERLINE).Println("root")
	out := filepath.Join(dest, "index.html")

	page := templates.RootPage{
		Background: root.Background,
		Logo:       root.Logo,
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
		page.Dates = fmt.Sprintf("%s - %s", startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
	}

	slices.SortFunc(page.Series, func(a, b templates.SeriesHeader) int {
		return a.Index - b.Index
	})

	if err := templates.RenderRootPage(out, page); err != nil {
		cmd.logError(err, "render ROOT failed")
		return errors.New("error rendering root")
	}

	style.Create.Printf("+ %s\n", out)
	cmd.logCreate(out)

	files := []string{"style.css", "script.js", root.Background, root.Logo}
	cmd.copyFiles(dest, data.PUBLIC_PATH, files)
	return nil
}

func (cmd BuildCommand) buildSubSeriesLinks(name string, series SeriesData) []templates.SubSeriesLink {
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

func (cmd BuildCommand) copyFiles(dest, src string, files []string) {
	stats := FileStats{}

	for _, f := range files {
		cmd.copyFileIfNewer(filepath.Join(dest, f), filepath.Join(src, f), &stats)
	}

	style.Create.Printf("[%d] files copied, ", stats.Copied)
	style.Info.Printf("[%d] files up-to-date, ", stats.UpToDate)
	style.Error.Printf("[%d] files not found, ", stats.NotFound)
	style.Error.Printf("[%d] files invalid\n", stats.Invalid)
}

func (cmd BuildCommand) copyFileIfNewer(dest, src string, stats *FileStats) {
	newer, err := cmd.isFileNewer(src, dest)
	if err != nil {
		cmd.logger.Printf("[NOT FOUND] %s\n", src)
		stats.NotFound++
		return
	}

	if !newer {
		cmd.logger.Printf("[UP_TO_DATE] %s\n", src)
		stats.UpToDate++
		return
	}

	style.Create.Printf("+ %s -> %s\n", src, dest)

	if err := file.Copy(dest, src); err != nil {
		cmd.logError(err, fmt.Sprintf("%s -> %s", src, dest))
		stats.Invalid++
	} else {
		cmd.logger.Printf("[COPIED] %s -> %s", src, dest)
		stats.Copied++
	}
}

func (cmd BuildCommand) isFileNewer(src, compare string) (bool, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return false, errors.Chain(err, "error reading source file")
	}
	srcTime := stat.ModTime().Truncate(time.Second)

	stat, err = os.Stat(compare)
	if err != nil {
		return true, nil
	}
	destTime := stat.ModTime().Truncate(time.Second)

	return destTime.Before(srcTime), nil
}

func (cmd BuildCommand) logError(err error, msg string) {
	cmd.logger.Printf("[ERROR] %s\n  %s\n", msg, err)
}

func (cmd BuildCommand) logCreate(file string) {
	cmd.logger.Printf("[CREATE] %s\n", file)
}

func (cmd BuildCommand) logBuild(series string) {
	cmd.logger.Printf("[BUILD] %s\n", series)
}

func (BuildCommand) printCreate(file string) {
	style.Create.Printf("[+] %s\n", file)
}

func (BuildCommand) formatDurationTimestamp(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDurationHMS(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%dh %dm %ds", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDateRange(startDate, endDate time.Time) string {
	if startDate.IsZero() && endDate.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s - %s", startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
}
