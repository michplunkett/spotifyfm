package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	"github.com/michplunkett/spotifyfm/util/constants"
	"github.com/michplunkett/spotifyfm/util/environment"
)

type SpotifyAuthHandler interface {
	Authenticate() *spotify.Client
}

type spotifyAuthHandler struct {
	auth   spotify.Authenticator
	ch     chan *spotify.Client
	config environment.EnvVars
	state  string
}

func newSpotifyAuthHandlerGeneric(permissions []string, e environment.EnvVars) SpotifyAuthHandler {
	state, _ := generateRandomString(23)
	return &spotifyAuthHandler{
		auth:   spotify.NewAuthenticator(constants.SpotifyRedirectURL, permissions...),
		ch:     make(chan *spotify.Client),
		config: e,
		state:  state,
	}
}

func NewSpotifyAuthHandlerAll(e environment.EnvVars) SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserFollowRead,
		spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail,
		spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistReadCollaborative,
		spotify.ScopeUserLibraryRead, spotify.ScopeUserTopRead,
		spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserReadCurrentlyPlaying)
	return newSpotifyAuthHandlerGeneric(permissions, e)
}

func NewSpotifyAuthHandlerProfile(e environment.EnvVars) SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserFollowRead,
		spotify.ScopeUserReadPrivate, spotify.ScopeUserReadEmail)
	return newSpotifyAuthHandlerGeneric(permissions, e)
}

func NewSpotifyAuthHandlerMusicStorage(e environment.EnvVars) SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryRead)
	return newSpotifyAuthHandlerGeneric(permissions, e)
}

func NewSpotifyAuthHandlerActivity(e environment.EnvVars) SpotifyAuthHandler {
	var permissions = make([]string, 0)
	permissions = append(permissions, spotify.ScopeUserTopRead, spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadCurrentlyPlaying)
	return newSpotifyAuthHandlerGeneric(permissions, e)
}

func (handler *spotifyAuthHandler) Authenticate() *spotify.Client {
	expirationTime := handler.config.GetSpotifyTokenExpiration()
	if !expirationTime.IsZero() && expirationTime.After(time.Now()) {
		token := &oauth2.Token{
			AccessToken:  handler.config.GetSpotifyToken(),
			TokenType:    "Bearer",
			RefreshToken: handler.config.GetSpotifyRefreshToken(),
			Expiry:       expirationTime,
		}
		client := handler.auth.NewClient(token)

		return &client
	} else {
		// first start an HTTP server
		http.HandleFunc("/spotify-callback", handler.finishAuthentication)

		authRequestUrl := handler.auth.AuthURL(handler.state)
		browser.OpenURL(authRequestUrl)

		// wait for auth to complete
		client := <-handler.ch

		return client
	}
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
	safeExpiration := tok.Expiry.Add(-time.Minute * 10)
	handler.config.SetSpotifyInfo(tok.AccessToken, tok.RefreshToken, safeExpiration)

	fmt.Fprintf(w, "Login Completed!")
	handler.ch <- &client
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
