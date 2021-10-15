package projects

import (
	"github.com/zmb3/spotify"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
)

type AudioFeatureProcessing interface {
	Execute()
}

type audioFeatureProcessing struct {
	audioFeatures  []*models.TrackAudioFeatures
	fileName       string
	spotifyHandler endpoints.SpotifyHandler
}

func NewAudioFeatureProcessing(fileName string, spotifyHandler endpoints.SpotifyHandler) AudioFeatureProcessing {
	return &audioFeatureProcessing{
		audioFeatures:  make([]*models.TrackAudioFeatures, 0),
		fileName:       fileName,
		spotifyHandler: spotifyHandler,
	}
}

func (a *audioFeatureProcessing) Execute() {
	a.parseInformationFromFile()
	a.getRecommendedTracks()
}

func (a *audioFeatureProcessing) parseInformationFromFile() {
	a.audioFeatures = append(a.audioFeatures, models.GetTrackAudioFeatures(a.fileName)...)
}

func (a *audioFeatureProcessing) getRecommendedTracks() {
	// Get stats for weekdays between 9-6
	work := a.getValuesBetweenTimeWeekdays(9, 18)
	// Get recommended tracks for mondayWork
	recs := a.spotifyHandler.GetTrackRecommendations(work, "hardcore", 30)
	recIds := make([]spotify.ID, 0)
	for _, track := range recs.Tracks {
		recIds = append(recIds, track.ID)
	}
	a.spotifyHandler.CreatePlaylistAndAddTracks("Work Hardcore Mix", "", recIds)
}

func (a *audioFeatureProcessing) getValuesBetweenTimeWeekdays(startHour, endHour int) *models.AttributeStats {
	acousticness := make([]float64, 0)
	danceability := make([]float64, 0)
	energy := make([]float64, 0)
	instrumentalness := make([]float64, 0)
	liveness := make([]float64, 0)
	loudness := make([]float64, 0)
	speechiness := make([]float64, 0)
	tempo := make([]float64, 0)
	valence := make([]float64, 0)

	for _, t := range a.audioFeatures {
		if int(t.ListenDate.Weekday()) != 0 && int(t.ListenDate.Weekday()) != 6 &&
			t.ListenDate.Hour() >= startHour && t.ListenDate.Hour() <= endHour {
			acousticness = append(acousticness, t.Acousticness)
			danceability = append(danceability, t.Danceability)
			energy = append(energy, t.Energy)
			instrumentalness = append(instrumentalness, t.Instrumentalness)
			liveness = append(liveness, t.Liveness)
			loudness = append(loudness, t.Loudness)
			speechiness = append(speechiness, t.Speechiness)
			tempo = append(tempo, t.Tempo)
			valence = append(valence, t.Valence)
		}
	}

	attr := &models.AttributeStats{}
	attr.Acousticness = models.BuildBasicStats(acousticness)
	attr.Danceability = models.BuildBasicStats(danceability)
	attr.Energy = models.BuildBasicStats(energy)
	attr.Instrumentalness = models.BuildBasicStats(instrumentalness)
	attr.Liveness = models.BuildBasicStats(liveness)
	attr.Loudness = models.BuildBasicStats(loudness)
	attr.Speechiness = models.BuildBasicStats(speechiness)
	attr.Tempo = models.BuildBasicStats(tempo)
	attr.Valence = models.BuildBasicStats(valence)

	return attr
}
