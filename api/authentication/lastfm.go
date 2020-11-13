package authentication

import (
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/browser"
	"github.com/shkh/lastfm-go/lastfm"

	"github.com/michplunkett/spotifyfm/config"
)

type LastFMAuthHandler interface {
	Authenticate() *lastfm.Api
	finishAuthentication(w http.ResponseWriter, r *http.Request)
}

var channelErr = make(chan error)

type lastFMAuthHandler struct {
	api    *lastfm.Api
	config config.EnvVars
}

func NewLastFMAuthHandler(e config.EnvVars) LastFMAuthHandler {
	return &lastFMAuthHandler{
		api:    lastfm.New(e.GetLastFMApiKey(), e.GetLastFMSharedSecret()),
		config: e,
	}
}

func (handler *lastFMAuthHandler) Authenticate() *lastfm.Api {
	// it's a non-empty string check because last.fm auth lasts FOREVER
	if handler.config.GetLastFMSessionKey() != config.EmptyString {
		handler.api.SetSession(handler.config.GetLastFMSessionKey())
	} else {
		http.HandleFunc("/lastfm-callback", handler.finishAuthentication)
		authRequestUrl := handler.api.GetAuthRequestUrl(config.LastFMRedirectURL)

		//fmt.Println("Opening the LastFM authorization URL in your browser:", authRequestUrl)
		browser.OpenURL(authRequestUrl)

		err := <-channelErr
		if err != nil {
			log.Fatal(err)
		}

		handler.config.SetLastFMSessionKey(handler.api.GetSessionKey())
	}

	return handler.api
}

func (handler *lastFMAuthHandler) finishAuthentication(w http.ResponseWriter, r *http.Request) {
	query, _ := url.ParseQuery(r.URL.RawQuery)
	token := query.Get("token")
	err := handler.api.LoginWithToken(token)
	channelErr <- err
}
