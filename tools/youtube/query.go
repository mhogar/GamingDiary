package youtube

import (
	"context"

	"github.com/binarysoupdev/go-extensions/errors"
	"google.golang.org/api/youtube/v3"
)

const MAX_VIDEO_QUERY = 50

func handleError(err error) error {
	return errors.Chain(err, "error sending response")
}

func (yt YTClient) GetMyChannel(parts ...string) (*youtube.Channel, error) {
	call := yt.service.Channels.List(parts).Mine(true)

	res, err := call.Do()
	if err != nil {
		return nil, handleError(err)
	}

	return res.Items[0], nil
}

func (yt YTClient) GetChannelByHandle(handle string, parts ...string) (*youtube.Channel, error) {
	call := yt.service.Channels.List(parts).ForHandle(handle)

	res, err := call.Do()
	if err != nil {
		return nil, handleError(err)
	}

	return res.Items[0], nil
}

func (yt YTClient) GetMyPlaylists(ctx context.Context, parts ...string) ([]*youtube.Playlist, error) {
	playlists := []*youtube.Playlist{}

	err := yt.service.Playlists.List(parts).Mine(true).Pages(ctx, func(res *youtube.PlaylistListResponse) error {
		playlists = append(playlists, res.Items...)
		return nil
	})

	if err != nil {
		return nil, handleError(err)
	}
	return playlists, nil
}

func (yt YTClient) GetPlaylists(c *youtube.Channel, parts ...string) ([]*youtube.Playlist, error) {
	call := yt.service.Playlists.List(parts).ChannelId(c.Id)

	res, err := call.Do()
	if err != nil {
		return nil, handleError(err)
	}

	return res.Items, nil
}

func (yt YTClient) GetItemsForPlaylist(ctx context.Context, id string) ([]*youtube.PlaylistItem, error) {
	items := []*youtube.PlaylistItem{}

	err := yt.service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(id).Pages(ctx, func(res *youtube.PlaylistItemListResponse) error {
		items = append(items, res.Items...)
		return nil
	})

	if err != nil {
		return nil, handleError(err)
	}
	return items, nil
}

func (yt YTClient) GetVideos(ctx context.Context, ids ...string) ([]*youtube.Video, error) {
	videos := []*youtube.Video{}

	for {
		min := min(len(ids), MAX_VIDEO_QUERY)

		err := yt.service.Videos.List([]string{"snippet", "fileDetails"}).Id(ids[0:min]...).Pages(ctx, func(res *youtube.VideoListResponse) error {
			videos = append(videos, res.Items...)
			return nil
		})

		if err != nil {
			return nil, handleError(err)
		}

		if len(ids) > MAX_VIDEO_QUERY {
			ids = ids[MAX_VIDEO_QUERY:]
		} else {
			break
		}
	}

	return videos, nil
}
