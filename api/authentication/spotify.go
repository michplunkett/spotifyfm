package authentication

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/config"
	"github.com/michplunkett/spotifyfm/utility"
)

type spotifyAPI struct {
	spotifyConfig *config.SpotifyConfig
}

type SpotifyAuthHandler interface {
	Authenticate() *spotify.Client
	finishAuthentication(w http.ResponseWriter, r *http.Request)
}

type spotifyAuthHandler struct {
	auth  spotify.Authenticator
	ch    chan *spotify.Client
	state string
}

func NewSpotifyAuthHandlerProfile() SpotifyAuthHandler {
	state, _ := utility.GenerateRandomString(23)
	return &spotifyAuthHandler{
		auth: spotify.NewAuthenticator(utility.RedirectURL, spotify.ScopeUserFollowRead,
			spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail),
		ch:    make(chan *spotify.Client),
		state: state,
	}
}

func NewSpotifyAuthHandlerMusicStorage() SpotifyAuthHandler {
	state, _ := utility.GenerateRandomString(23)
	return &spotifyAuthHandler{
		auth: spotify.NewAuthenticator(utility.RedirectURL, spotify.ScopePlaylistReadPrivate,
			spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryRead),
		ch:    make(chan *spotify.Client),
		state: state,
	}
}

func NewSpotifyAuthHandlerActivity() SpotifyAuthHandler {
	state, _ := utility.GenerateRandomString(23)
	return &spotifyAuthHandler{
		auth: spotify.NewAuthenticator(utility.RedirectURL, spotify.ScopeUserTopRead, spotify.ScopeUserReadRecentlyPlayed,
			spotify.ScopeUserReadCurrentlyPlaying),
		ch:    make(chan *spotify.Client),
		state: state,
	}
}

func (handler *spotifyAuthHandler) Authenticate() *spotify.Client {
	// first start an HTTP server
	http.HandleFunc("/callback", handler.finishAuthentication)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	authRequestUrl := handler.auth.AuthURL(handler.state)
	fmt.Println("Opening the Spotify authorization URL in your browser:", authRequestUrl)
	browser.OpenURL(authRequestUrl)

	// wait for auth to complete
	client := <-handler.ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.Followers)

	return client
}

func (handler *spotifyAuthHandler) finishAuthentication(w http.ResponseWriter, r *http.Request) {
	tok, err := handler.auth.Token(handler.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != handler.state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, handler.state)
	}
	// use the token to get an authenticated client
	client := handler.auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	handler.ch <- &client
}