package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	redirectPort     = 8081
	callbackEndpoint = "/callback"
	state            = "abc123"
	tokenFile        = "token.json"
	credsFile        = "credentials.json"
)

var (
	redirectURI = fmt.Sprintf("http://localhost:%d%s", redirectPort, callbackEndpoint)
)

func getClient() *spotify.Client {
	ctx := context.Background()

	ch := make(chan *oauth2.Token)
	var tok *oauth2.Token

	if f, err := os.ReadFile(tokenFile); err != nil {
		credRaw, err := ioutil.ReadFile("credentials.json")
		if err != nil {
			panic(err)
		}
		creds := make(map[string]string)
		err = json.Unmarshal(credRaw, &creds)
		if err != nil {
			panic(err)
		}

		// Create our auth object
		auth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopePlaylistReadPrivate), spotifyauth.WithClientID(creds["id"]), spotifyauth.WithClientSecret(creds["secret"]))

		// first start an HTTP server
		http.HandleFunc(callbackEndpoint, func(w http.ResponseWriter, r *http.Request) {
			_tok, err := auth.Token(r.Context(), state, r)
			if err != nil {
				http.Error(w, "Couldn't get token", http.StatusForbidden)
				log.Fatal(err)
			}
			if st := r.FormValue("state"); st != state {
				http.NotFound(w, r)
				log.Fatalf("State mismatch: %s != %s\n", st, state)
			}

			tokStr, err := json.Marshal(_tok)
			if err != nil {
				http.Error(w, "Couldn't marshal token", http.StatusInternalServerError)
				log.Fatal(err)
			}
			err = os.WriteFile(tokenFile, tokStr, 0666)
			if err != nil {
				http.Error(w, "Couldn't write token", http.StatusInternalServerError)
				log.Fatal(err)
			}

			fmt.Fprintf(w, "Login Completed!")

			ch <- _tok
		})
		go func() {
			err := http.ListenAndServe(":"+fmt.Sprint(redirectPort), nil)
			if err != nil {
				log.Fatal(err)
			}
		}()

		url := auth.AuthURL(state)
		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		tok = <-ch
	} else {
		tok = new(oauth2.Token)

		err := json.Unmarshal(f, &tok)
		if err != nil {
			log.Fatal(err)
		}
	}

	// use the token to get an authenticated client
	httpClient := spotifyauth.New().Client(ctx, tok)
	client := spotify.New(httpClient)
	if client == nil {
		log.Fatalf("client is nil")
	}

	return client
}
