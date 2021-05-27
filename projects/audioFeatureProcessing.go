package projects

import (
	"fmt"

	"github.com/michplunkett/spotifyfm/models"
)

type AudioFeatureProcessing interface {
	Execute()
}

type audioFeatureProcessing struct {
	audioFeatures []*models.TrackAudioFeatures
	fileName      string
}

func NewAudioFeatureProcessing(fileName string) AudioFeatureProcessing {
	return &audioFeatureProcessing{
		audioFeatures: make([]*models.TrackAudioFeatures, 0),
		fileName:      fileName,
	}
}

func (a *audioFeatureProcessing) Execute() {
	a.parseInformationFromFile()
}

func (a *audioFeatureProcessing) parseInformationFromFile() {
	a.audioFeatures = append(a.audioFeatures, models.GetTrackAudioFeatures(a.fileName)...)
	// Get stats for Monday between 9-6
	mondayWork := a.getValuesBetweenDayAndTime(1, 9, 18)
	fmt.Println(mondayWork.Acousticness.Median)
}

func (a *audioFeatureProcessing) getValuesBetweenDayAndTime(day, startHour, endHour int) *models.AttributeStats {
	acousticness := make([]float32, 0)
	danceability := make([]float32, 0)
	energy := make([]float32, 0)
	instrumentalness := make([]float32, 0)
	liveness := make([]float32, 0)
	loudness := make([]float32, 0)
	speechiness := make([]float32, 0)
	tempo := make([]float32, 0)
	valence := make([]float32, 0)

	for _, t := range a.audioFeatures {
		if int(t.ListenDate.Weekday()) == day && t.ListenDate.Hour() >= startHour && t.ListenDate.Hour() <= endHour {
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
