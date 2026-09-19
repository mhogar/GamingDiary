package cmds

import (
	"app/data"
	"app/data/series"
	youtube_data "app/data/youtube"
	"app/tools/youtube"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type YoutubeCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewYoutubeCommand() *YoutubeCommand {
	return &YoutubeCommand{
		CommandBase: command.NewCommandBase("youtube", "Build entries from Youtube data"),
	}
}

func (cmd *YoutubeCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd YoutubeCommand) Run(args []string) error {
	s := cmd.Flags.String("series", "", "name of the series")
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	cmd.Flags.Parse(args)

	if *s == "" {
		return errors.New("\"series\" cannot be empty")
	}

	series, err := series.Select(*s)
	if err != nil {
		return err
	}
	style.BoldInfo.Println(*s)

	return cmd.downloadData(series, *forceAuth)
}

func (cmd YoutubeCommand) downloadData(series series.Series, forceAuth bool) error {
	meta, err := json.UnmarshalFile[youtube_data.Meta](filepath.Join(data.STATIC_PATH, series.GetName(), "youtube.json"))
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
