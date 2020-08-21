package authentication

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/browser"
	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/config"
	"github.com/michplunkett/spotifyfm/utility"
)

type LastFMAuthHandler interface {
	Authenticate() *lastfm.Api
	finishAuthentication(w http.ResponseWriter, r *http.Request)
}

var channelErr = make(chan error)

type lastFMAuthHandler struct {
	api    *lastfm.Api
	key    string
	secret string
}

func NewLastFMAuthHandler(config *config.LastFmConfig) LastFMAuthHandler {
	return &lastFMAuthHandler{
		api: lastfm.New(config.GetApiKey(), config.GetSharedSecret()),
	}
}

func (handler *lastFMAuthHandler) Authenticate() *lastfm.Api {
	http.HandleFunc("/callback", handler.finishAuthentication)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)
	authRequestUrl := handler.api.GetAuthRequestUrl(utility.RedirectURL)

	fmt.Println("Opening the LastFM authorization URL in your browser:", authRequestUrl)
	browser.OpenURL(authRequestUrl)

	err := <-channelErr
	if err != nil {
		log.Fatal(err)
	}

	result, _ := handler.api.User.GetInfo(nil)
	fmt.Println("We have a login!: ", result)

	return handler.api
}

func (handler *lastFMAuthHandler) finishAuthentication(w http.ResponseWriter, r *http.Request) {
	query, _ := url.ParseQuery(r.URL.RawQuery)
	token := query.Get("token")
	err := handler.api.LoginWithToken(token)
	channelErr <- err
}
