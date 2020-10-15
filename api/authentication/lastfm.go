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
	http.HandleFunc("/lastfm-callback", handler.finishAuthentication)
	authRequestUrl := handler.api.GetAuthRequestUrl(utility.LastFMRedirectURL)

	//fmt.Println("Opening the LastFM authorization URL in your browser:", authRequestUrl)
	browser.OpenURL(authRequestUrl)

	err := <-channelErr
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("---------------------")
	fmt.Println("LAST.FM THANGS")

	userInfo, _ := handler.api.User.GetInfo(nil)
	fmt.Println("You are logged in as ", userInfo.RealName)

	topTracksParam := make(map[string]interface{})
	topTracksParam["user"] = userInfo.Name
	topTracksParam["limit"] = 50
	topTracksParam["period"] = "1month"
	topTracks, _ := handler.api.User.GetTopTracks(topTracksParam)
	fmt.Println("This is your top track ", topTracks.Tracks[0].Name)

	return handler.api
}

func (handler *lastFMAuthHandler) finishAuthentication(w http.ResponseWriter, r *http.Request) {
	query, _ := url.ParseQuery(r.URL.RawQuery)
	token := query.Get("token")
	err := handler.api.LoginWithToken(token)
	channelErr <- err
}
