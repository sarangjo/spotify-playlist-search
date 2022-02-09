package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
)

func main() {
	client := getClient()

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	fmt.Println("Your display name is:", user.DisplayName)

	playlists, err := client.CurrentUsersPlaylists(context.Background(), spotify.Limit(50))
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range playlists.Playlists {
		fmt.Println("Playlist:\t", p.Name)
	}
}
