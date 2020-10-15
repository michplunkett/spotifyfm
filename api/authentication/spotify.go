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

func newSpotifyAuthHandlerGeneric(permissions []string) SpotifyAuthHandler {
	state, _ := utility.GenerateRandomString(23)
	return &spotifyAuthHandler{
		auth: spotify.NewAuthenticator(utility.SpotifyRedirectURL, permissions...),
		ch:    make(chan *spotify.Client),
		state: state,
	}
}

func NewSpotifyAuthHandlerAll() SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserFollowRead,
		spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail,
		spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistReadCollaborative,
		spotify.ScopeUserLibraryRead, spotify.ScopeUserTopRead,
		spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadCurrentlyPlaying)
	return newSpotifyAuthHandlerGeneric(permissions)
}

func NewSpotifyAuthHandlerProfile() SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserFollowRead,
		spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail)
	return newSpotifyAuthHandlerGeneric(permissions)
}

func NewSpotifyAuthHandlerMusicStorage() SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryRead)
	return newSpotifyAuthHandlerGeneric(permissions)
}

func NewSpotifyAuthHandlerActivity() SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserTopRead, spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadCurrentlyPlaying)
	return newSpotifyAuthHandlerGeneric(permissions)
}

func (handler *spotifyAuthHandler) Authenticate() *spotify.Client {
	// first start an HTTP server
	http.HandleFunc("/spotify-callback", handler.finishAuthentication)

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
	fmt.Println("You are logged in as:", user.DisplayName)

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
