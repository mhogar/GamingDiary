package cmds

import (
	"context"
	"fmt"
	"gamingdiary/data"
	"gamingdiary/data/series"
	"gamingdiary/tools/youtube"
	"path/filepath"
	"regexp"
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
}

func (cmd *YoutubeCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd YoutubeCommand) Run(args []string) error {
	s := cmd.Flags.String("series", "", "name of the series")
	cache := cmd.Flags.String("cache", "", "use an existing cached data")
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	cmd.Flags.Parse(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}
	style.BoldInfo.Println(*s)

	series, err := series.Select(*s)
	if err != nil {
		return err
	}

	var videos []*youtube.Video
	if *cache != "" {
		videos, err = cmd.loadCachedData(*cache)
	} else {
		videos, err = cmd.downloadNewData(series, *forceAuth)
	}
	if err != nil {
		return err
	}

	if err := cmd.createEntries(series, videos); err != nil {
		return errors.Chain(err, "error creating entires")
	}
	return nil
}

func (cmd YoutubeCommand) createEntries(series series.Series, videos []*youtube.Video) error {
	var errs errors.Errors

	// TODO: sort by publish date
	// TODO: better index system?

	for i, video := range videos {
		index := fmt.Sprintf("%02d", i)
		path := filepath.Join(data.STATIC_DIR, series.GetName(), fmt.Sprintf("entry%s.json", index))

		if err := cmd.createEntry(path, index, series, video); err != nil {
			errs.Add(errors.Format("[%s] %s", index, err))
		}
	}
	fmt.Println()

	return errs.Collapse("\n  ")
}

func (cmd YoutubeCommand) createEntry(path, index string, series series.Series, video *youtube.Video) error {
	duration, err := cmd.parseDuration(video.ContentDetails.Duration)
	if err != nil {
		return err
	}

	date, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
	if err != nil {
		return errors.Chain(err, "error parsing date")
	}

	entry := data.Entry{
		YoutubeId: video.Id,
		Title:     video.Snippet.Title,
		Date:      date.Format(data.ENTRY_DATE_FORMAT),
		Duration:  float32(duration),
		Thumbnail: fmt.Sprintf("t%s.png", index),
		Video:     fmt.Sprintf("v%s.mp4", index),
	}

	if err := series.BuildEntryFromYoutube(index, video, &entry); err != nil {
		return err
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}

	fmt.Printf("\r... %s ", style.Create.Sprintf("+ %s ", path))
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

func (cmd YoutubeCommand) downloadNewData(series series.Series, forceAuth bool) ([]*youtube.Video, error) {
	meta, err := json.UnmarshalFile[data.YoutubeMeta](filepath.Join("series", series.GetName(), "youtube.json"))
	if err != nil {
		return nil, errors.Chain(err, "error reading youtube meta file")
	}

	ctx := context.Background()

	client, err := youtube.NewClient(ctx, forceAuth)
	if err != nil {
		return nil, errors.Chain(err, "error creating youtube client")
	}

	ids, err := cmd.loadVideoIdsFromPlaylist(client, ctx, meta.Playlist)
	if err != nil {
		return nil, err
	}

	videos, err := cmd.downloadVideoData(client, ctx, ids)
	if err != nil {
		return nil, err
	}

	output := fmt.Sprintf("youtube/%s_%s.json", strings.ReplaceAll(series.GetName(), "/", "_"), time.Now().Format("2006-01-02_15:04:05"))
	if err := json.MarshalFilePretty(videos, output, "  "); err != nil {
		return nil, errors.Chain(err, "error saving videos json")
	}

	style.Create.Printf("+ %s\n", output)
	return videos, nil
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
