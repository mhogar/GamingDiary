package oauth

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/got-style/style"
	"golang.org/x/oauth2"
)

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	style.Bold.Println("Authorizing from Web")

	server := newTokenServer()
	server.Start()

	config.RedirectURL = server.URL()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	err := launchAuthURL(authURL)
	if err != nil {
		return nil, errors.Chain(err, "error launching auth URL")
	}
	fmt.Println("  opened request in browser")

	code := <-server.Code
	fmt.Println("  received code")

	server.Close()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.Chain(err, "error retrieving token from web")
	}

	return token, nil
}

func launchAuthURL(url string) error {
	err := exec.Command("xdg-open", url).Start()
	if err != nil {
		return errors.Chain(err, "error launching browser")
	}
	return nil
}
