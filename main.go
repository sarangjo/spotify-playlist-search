package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var c chan string

const (
	tokenFile   = "token.json"
	redirectUri = "http://localhost:8081/callback"
)

func main() {
	ctx := context.Background()
	credRaw, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}
	creds := make(map[string]string)
	err = json.Unmarshal(credRaw, &creds)
	if err != nil {
		panic(err)
	}

	config := &clientcredentials.Config{
		ClientID:     creds["id"],
		ClientSecret: creds["secret"],
		TokenURL:     spotifyauth.TokenURL,
		Scopes:       []string{"playlist-read-private"},
	}

	token, err := config.Token(ctx)

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your id is:", user.ID)
	fmt.Println("Your display name is:", user.DisplayName)

	playl, err := client.GetPlaylistsForUser(ctx, user.ID, spotify.Limit(50))
	if err != nil {
		panic(err)
	}

	for _, p := range playl.Playlists {
		if p.Owner.ID != user.ID {
			continue
		}

		ep := p.Tracks.
	}
}
