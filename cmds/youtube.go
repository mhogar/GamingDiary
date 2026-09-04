package cmds

import (
	"context"
	"fmt"
	"gamingdiary/data"
	client "gamingdiary/tools/youtube"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
	"google.golang.org/api/youtube/v3"
)

func NewYouTubeCommand() *YouTubeCommand {
	return &YouTubeCommand{
		CommandBase: command.NewCommandBase("youtube", "download data from YouTube"),
	}
}

type YouTubeCommand struct {
	command.CommandBase
	command.FlagCommand

	iterator       int
	iteratorFormat string
}

func (cmd *YouTubeCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd YouTubeCommand) Run(args []string) error {
	forceAuth := cmd.Flags.Bool("auth", false, "force re-authentication")
	download := cmd.Flags.String("download", "download", "the download directory")
	name := cmd.Flags.String("name", "", "the name of the series")
	num := cmd.Flags.String("num", "0", "iterator starting value and padding")
	cmd.Flags.Parse(args)

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}
	dataPath := filepath.Join("data", *name)

	iter64, err := strconv.ParseInt(*num, 10, 16)
	if err != nil {
		return errors.Chain(err, "invalid iterator")
	}
	cmd.iterator = int(iter64)
	cmd.iteratorFormat = fmt.Sprintf("%%0%dd", len(*num))

	series, err := json.UnmarshalFile[data.Series](filepath.Join(dataPath, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading index file")
	}
	style.BoldInfo.Println(*name)

	ctx := context.Background()

	yt, err := client.NewClient(ctx, *forceAuth)
	if err != nil {
		return errors.Chain(err, "error creating youtube client")
	}

	ids, err := cmd.loadVideoIdsPlaylist(yt, ctx, series.YoutubePlaylist)
	if err != nil {
		return err
	}

	output := filepath.Join(dataPath, *download, time.Now().Format(time.DateTime))
	if err := os.MkdirAll(output, 0755); err != nil {
		return errors.Chain(err, "error creating download directory")
	}

	return cmd.downloadVideoMeta(output, yt, ctx, ids)
}

func (cmd YouTubeCommand) loadVideoIdsPlaylist(yt *client.YTClient, ctx context.Context, playlist string) ([]string, error) {
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

func (cmd YouTubeCommand) downloadVideoMeta(path string, client *client.YTClient, ctx context.Context, ids []string) error {
	fmt.Print("Downloading video metadata")

	videos, err := client.GetVideos(ctx, ids...)
	if err != nil {
		return errors.Chain(err, "error getting videos")
	}
	fmt.Printf(" -> %s\n", style.BoldInfo.Sprintf("[%d] downloaded", len(videos)))

	for _, video := range videos {
		name := fmt.Sprintf("meta%s.txt", fmt.Sprintf(cmd.iteratorFormat, cmd.iterator))

		if err := cmd.saveMeta(filepath.Join(path, name), video); err != nil {
			return err
		}
		cmd.iterator++
	}
	fmt.Println()

	return nil
}

func (cmd YouTubeCommand) saveMeta(path string, v *youtube.Video) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating meta file")
	}
	defer f.Close()

	style.Create.Printf("\r%s -> %s", v.Id, path)

	fmt.Fprintf(f, "https://www.youtube.com/watch?v=%s\n\n%s\n---\n\n%s\n", v.Id, v.Snippet.Title, v.Snippet.Description)
	return nil
}
