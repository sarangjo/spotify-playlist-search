package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var c chan string

const (
	redirectUri = "http://localhost:8081/callback"
)

func main() {
	c = make(chan string)

	go server2()

	credRaw, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}
	creds := make(map[string]string)
	err = json.Unmarshal(credRaw, &creds)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     creds["id"],
		ClientSecret: creds["secret"],
		Scopes:       []string{"playlist-read-private"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
		RedirectURL: redirectUri,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	// Wait on code
	code := <-c

	// Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	// Handoff client
	client := conf.Client(ctx, tok)
	process(client)
}

func process(client *http.Client) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	res := make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	fmt.Println("Your user id is", res["id"].(string))
	fmt.Println("Your display name is", res["display_name"].(string))
}

func server2() {
	http.HandleFunc("/callback", callback2)

	http.ListenAndServe(":8081", nil)
}

func callback2(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	code := q.Get("code")

	c <- code

	w.Write([]byte("success"))
}
