package cmds

import (
	"context"
	"fmt"
	"gamingdiary/data"
	"gamingdiary/data/series"
	"gamingdiary/tools/youtube"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

var YT_DURATION_REGEX = regexp.MustCompile(`([0-9]+)([^0-9])`)

func NewYoutubeCommand() *YoutubeCommand {
	return &YoutubeCommand{
		CommandBase: command.NewCommandBase("youtube", "Build entries from Youtube data"),
	}
}

type YoutubeCommand struct {
	command.CommandBase
	command.FlagCommand

	indexStart  int
	indexFormat string
}

func (cmd *YoutubeCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd YoutubeCommand) Run(args []string) error {
	s := cmd.Flags.String("series", "", "name of the series")
	parse := cmd.Flags.String("parse", "", "parse existing data")
	download := cmd.Flags.Bool("download", false, "download new data")
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	index := cmd.Flags.String("index", "00", "starting index and padding")
	cmd.Flags.Parse(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}
	if !*download && *parse == "" {
		return errors.New("\"parse\" cannot be empty")
	}

	series, err := series.Select(*s)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(*s)

	i64, err := strconv.ParseInt(*index, 10, 16)
	if err != nil {
		return errors.Chain(err, "invalid index")
	}
	cmd.indexStart = int(i64)
	cmd.indexFormat = fmt.Sprintf("%%0%dd", len(*index))

	if *download {
		return cmd.downloadData(series, *forceAuth)
	}

	videos, err := cmd.loadCachedData(*parse)
	if err != nil {
		return err
	}
	return cmd.createEntries(series, videos)
}

func (cmd YoutubeCommand) createEntries(series series.Series, videos []*youtube.Video) error {
	type video struct {
		Video *youtube.Video
		Date  time.Time
	}

	videosByDate := make([]video, len(videos))
	for i, v := range videos {
		date, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)
		if err != nil {
			return errors.Chain(err, "error parsing date")
		}

		videosByDate[i] = video{
			Video: v,
			Date:  date,
		}
	}
	slices.SortFunc(videosByDate, func(a, b video) int {
		return a.Date.Compare(b.Date)
	})

	for i, v := range videosByDate {
		index := i + cmd.indexStart
		path := filepath.Join(data.STATIC_DIR, series.GetName(), fmt.Sprintf("entry%s.json", fmt.Sprintf(cmd.indexFormat, index)))

		err := cmd.createEntry(path, index, series, v.Video, v.Date)
		if err == nil {
			fmt.Printf("\r... %s ", style.Create.Sprintf("[+] %s ", path))
		} else {
			style.Error.Printf("\n[x] %s\n", err)
		}
	}
	fmt.Println()

	return nil
}

func (cmd YoutubeCommand) createEntry(path string, index int, series series.Series, video *youtube.Video, date time.Time) error {
	duration, err := cmd.parseDuration(video.ContentDetails.Duration)
	if err != nil {
		return err
	}

	entry := data.Entry{
		Title:        video.Snippet.Title,
		Date:         date,
		Duration:     float32(duration),
		Thumbnail:    fmt.Sprintf("t%s.png", fmt.Sprintf(cmd.indexFormat, index)),
		Video:        fmt.Sprintf("v%s.mp4", fmt.Sprintf(cmd.indexFormat, index)),
		YoutubeId:    video.Id,
		YoutubeVideo: fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id),
	}

	if video.Snippet.Thumbnails.Maxres != nil {
		entry.YoutubeThumbnail = video.Snippet.Thumbnails.Maxres.Url
	} else {
		entry.YoutubeThumbnail = video.Snippet.Thumbnails.Medium.Url
	}

	if err := series.BuildEntryFromYoutube(index, video, &entry); err != nil {
		return err
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}
	return nil
}

func (cmd YoutubeCommand) parseDuration(str string) (int64, error) {
	matches := YT_DURATION_REGEX.FindAllStringSubmatch(str, 2)
	if len(matches) == 0 {
		return 0, errors.Format("invalid duration format \"%s\"", str)
	}

	var duration int64
	for _, match := range matches {
		d, _ := strconv.ParseInt(match[1], 10, 16)

		if match[2] == "M" {
			duration += d * 60
		} else {
			duration += d
		}
	}
	return duration, nil
}

func (cmd YoutubeCommand) loadCachedData(path string) ([]*youtube.Video, error) {
	videos, err := json.UnmarshalFile[[]*youtube.Video](path)
	if err != nil {
		return nil, errors.Chain(err, "error loading youtube video cache")
	}
	return videos, nil
}

func (cmd YoutubeCommand) downloadData(series series.Series, forceAuth bool) error {
	meta, err := json.UnmarshalFile[data.YoutubeMeta](filepath.Join(data.STATIC_DIR, series.GetName(), "youtube.json"))
	if err != nil {
		return errors.Chain(err, "error reading youtube meta file")
	}

	ctx := context.Background()

	client, err := youtube.NewClient(ctx, forceAuth)
	if err != nil {
		return errors.Chain(err, "error creating youtube client")
	}

	ids, err := cmd.loadVideoIdsFromPlaylist(client, ctx, meta.Playlist)
	if err != nil {
		return err
	}

	videos, err := cmd.downloadVideoData(client, ctx, ids)
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s_%s.json", strings.ReplaceAll(series.GetName(), "/", "_"), time.Now().Format("2006-01-02_15:04:05"))
	output := filepath.Join(data.YOUTUBE_DATA_PATH, file)

	if err := json.MarshalFilePretty(videos, output, "  "); err != nil {
		return errors.Chain(err, "error saving videos json")
	}

	style.Create.Printf("+ %s\n", output)
	return nil
}

func (cmd YoutubeCommand) loadVideoIdsFromPlaylist(yt *youtube.YTClient, ctx context.Context, playlist string) ([]string, error) {
	if playlist == "" {
		return nil, errors.New("playlist ID cannot be empty")
	}

	fmt.Printf("Finding videos for YouTube playlist %s", style.Bold.Sprint(playlist))

	items, err := yt.GetItemsForPlaylist(ctx, playlist)
	if err != nil {
		return nil, errors.Chain(err, "error getting playlist items")
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ContentDetails.VideoId
	}

	fmt.Print(" -> ")
	style.BoldInfo.Printf("[%d] videos found\n", len(ids))

	return ids, nil
}

func (cmd YoutubeCommand) downloadVideoData(client *youtube.YTClient, ctx context.Context, ids []string) ([]*youtube.Video, error) {
	fmt.Print("Downloading video data")

	videos, err := client.GetVideos(ctx, []string{"snippet", "contentDetails"}, ids...)
	if err != nil {
		return nil, errors.Chain(err, "error getting videos")
	}
	fmt.Printf(" -> %s\n", style.BoldInfo.Sprintf("[%d] downloaded", len(videos)))

	return videos, nil
}
