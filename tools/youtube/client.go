package youtube

import (
	"context"
	"local/tools/oauth"
	"os"

	"github.com/binarysoupdev/go-extensions/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const SECRET_FILE = ".secrets/youtube.json"

type YTClient struct {
	service *youtube.Service
}

func NewClient(ctx context.Context, forceAuth bool) (*YTClient, error) {
	bytes, err := os.ReadFile(SECRET_FILE)
	if err != nil {
		return nil, errors.Chain(err, "error reading secret file")
	}

	// NOTE: delete cached token if changing scopes
	cfg, err := google.ConfigFromJSON(bytes, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, errors.Chain(err, "error creating config from client secret")
	}

	client, err := oauth.NewClient(ctx, cfg, forceAuth)
	if err != nil {
		return nil, errors.Chain(err, "error creating client")
	}

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Chain(err, "error creating service")
	}

	return &YTClient{
		service: service,
	}, nil
}
