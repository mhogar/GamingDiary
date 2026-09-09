package cmds

import (
	"context"
	"fmt"
	yt_data "gamingdiary/data/youtube"
	"gamingdiary/tools/youtube"
	"path/filepath"
	"strings"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func NewYouTubeCommand() *YouTubeCommand {
	return &YouTubeCommand{
		CommandBase: command.NewCommandBase("youtube", "Download data from Youtube"),
	}
}

type YouTubeCommand struct {
	command.CommandBase
	command.FlagCommand
}

func (cmd *YouTubeCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd YouTubeCommand) Run(args []string) error {
	series := cmd.Flags.String("series", "", "the name of the series")
	out := cmd.Flags.String("out", "youtube", "the output directory")
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	cmd.Flags.Parse(args)

	if *series == "" {
		return errors.New("\"series\" cannot be empty")
	}

	meta, err := json.UnmarshalFile[yt_data.Meta](filepath.Join("series", *series, "youtube.json"))
	if err != nil {
		return errors.Chain(err, "error reading youtube meta file")
	}
	style.BoldInfo.Println(*series)

	ctx := context.Background()

	yt, err := youtube.NewClient(ctx, *forceAuth)
	if err != nil {
		return errors.Chain(err, "error creating youtube client")
	}

	ids, err := cmd.loadVideoIdsFromPlaylist(yt, ctx, meta.Playlist)
	if err != nil {
		return err
	}

	videos, err := cmd.downloadVideoData(yt, ctx, ids)
	if err != nil {
		return err
	}

	output := filepath.Join(*out, fmt.Sprintf("%s_%s.json", strings.ReplaceAll(*series, "/", "_"), time.Now().Format("2006-01-02_15:04:05")))
	if err := json.MarshalFilePretty(videos, output, "  "); err != nil {
		return errors.Chain(err, "error saving videos json")
	}

	style.Create.Printf("+ %s\n", output)
	return nil
}

func (cmd YouTubeCommand) loadVideoIdsFromPlaylist(yt *youtube.YTClient, ctx context.Context, playlist string) ([]string, error) {
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

func (cmd YouTubeCommand) downloadVideoData(client *youtube.YTClient, ctx context.Context, ids []string) ([]*youtube.Video, error) {
	fmt.Print("Downloading video data")

	videos, err := client.GetVideos(ctx, ids...)
	if err != nil {
		return nil, errors.Chain(err, "error getting videos")
	}
	fmt.Printf(" -> %s\n", style.BoldInfo.Sprintf("[%d] downloaded", len(videos)))

	return videos, nil
}
