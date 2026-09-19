package entry_cmd

import (
	"app/data"
	"app/data/build"
	"app/data/series"
	"app/tools/youtube"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

func (cmd EntryCommand) createEntriesFromYoutube(cacheFile string, series series.Series) error {
	type video struct {
		Video *youtube.Video
		Date  time.Time
	}

	videos, err := json.UnmarshalFile[[]*youtube.Video](cacheFile)
	if err != nil {
		return errors.Chain(err, "error loading youtube video cache")
	}

	videosByDate := make([]video, len(videos))
	for i, v := range videos {
		date, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)
		if err != nil {
			return errors.Chain(err, "error parsing youtube date")
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
		index, path := cmd.calcNextEntryFromIndex(series.GetName(), i)

		if err := cmd.createYoutubeEntry(path, index, i, series, v.Video, v.Date); err != nil {
			style.Error.Printf("x %s\n", err)
		}
	}
	return nil
}

func (cmd EntryCommand) createYoutubeEntry(path string, index, videoIndex int, series series.Series, video *youtube.Video, date time.Time) error {
	duration, err := cmd.parseYoutubeDuration(video.ContentDetails.Duration)
	if err != nil {
		return err
	}

	entry := build.Entry{
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

	if err := series.BuildEntryFromYoutube(index, videoIndex, video, &entry); err != nil {
		return err
	}

	if err := json.MarshalFilePretty(entry, path, "    "); err != nil {
		return errors.Chain(err, "error saving entry file")
	}

	style.Create.Printf("+ %s\n", path)
	return nil
}

func (cmd EntryCommand) parseYoutubeDuration(str string) (int64, error) {
	matches := data.YOUTUBE_DURATION_REGEX.FindAllStringSubmatch(str, 2)
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
