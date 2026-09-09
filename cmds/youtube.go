package cmds

import (
	"context"
	"fmt"
	"gamingdiary/data"
	yt_data "gamingdiary/data/youtube"
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
		CommandBase: command.NewCommandBase("youtube", "Download data from Youtube"),
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
	series := cmd.Flags.String("series", "", "the name of the series")
	cache := cmd.Flags.String("cache", "", "use an existing cached data")
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	cmd.Flags.Parse(args)

	if *series == "" {
		return errors.New("\"series\" cannot be empty")
	}
	style.BoldInfo.Println(*series)

	var videos []*youtube.Video
	var err error

	if *cache != "" {
		videos, err = cmd.loadCachedData(*cache)
	} else {
		videos, err = cmd.downloadNewData(*series, *forceAuth)
	}
	if err != nil {
		return err
	}

	if err := cmd.createEntries(filepath.Join("series", *series), videos); err != nil {
		return errors.Chain(err, "error creating entires")
	}
	return nil
}

func (cmd YoutubeCommand) createEntries(path string, videos []*youtube.Video) error {
	var errs errors.Errors

	// TODO: sort by publish date
	// TODO: better index system?

	for i, video := range videos {
		index := fmt.Sprintf("%02d", i)

		if err := cmd.createEntry(filepath.Join(path, fmt.Sprintf("entry%s.json", index)), video); err != nil {
			errs.Add(errors.Format("[%s] %s", index, err))
		}
	}

	return errs.Collapse("\n  ")
}

func (cmd YoutubeCommand) createEntry(path string, video *youtube.Video) error {
	duration, err := cmd.parseDuration(video.ContentDetails.Duration)
	if err != nil {
		return err
	}

	entry := data.Entry{
		Date:     video.Snippet.PublishedAt,
		Duration: float32(duration),
		Youtube:  video.Id,
	}

	//TODO: parse from series

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}

	style.Create.Printf("+ %s\n", path)
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

func (cmd YoutubeCommand) downloadNewData(series string, forceAuth bool) ([]*youtube.Video, error) {
	meta, err := json.UnmarshalFile[yt_data.Meta](filepath.Join("series", series, "youtube.json"))
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

	output := fmt.Sprintf("youtube/%s_%s.json", strings.ReplaceAll(series, "/", "_"), time.Now().Format("2006-01-02_15:04:05"))
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
