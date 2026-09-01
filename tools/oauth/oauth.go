package oauth

import (
	"context"
	"log"
	"net/http"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"

	"golang.org/x/oauth2"
)

const TOKEN_FILE = ".secrets/auth_token.json"

func NewClient(ctx context.Context, cfg *oauth2.Config, forceAuth bool) (*http.Client, error) {
	token, err := loadToken(cfg, forceAuth)
	if err != nil {
		return nil, errors.Chain(err, "error getting token")
	}

	err = json.MarshalFile(token, TOKEN_FILE)
	if err != nil {
		log.Println("error saving token", err)
	}

	return cfg.Client(ctx, token), nil
}

func loadToken(cfg *oauth2.Config, forceAuth bool) (*oauth2.Token, error) {
	if forceAuth {
		style.BoldInfo.Println("[FORCED RE-AUTH]")
		return getTokenFromWeb(cfg)
	}

	token, err := json.UnmarshalFile[oauth2.Token](TOKEN_FILE)
	if err == nil {
		style.BoldInfo.Println("[AUTH CACHED]")
		return &token, nil
	}

	style.Error.Println("[NO AUTH CACHE]")
	return getTokenFromWeb(cfg)
}
