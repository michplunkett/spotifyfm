package cmd

import (
	"fmt"

	"github.com/michplunkett/spotifyfm/api/endpoints"
	"github.com/michplunkett/spotifyfm/models"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

func NewSpotifymd() *cobra.Command {
	return &cobra.Command{
		Use:   "spotify",
		Short: "Runs Spotify subcommands.",
	}
}

func NewAudioFeatureProcessing(spotifyHandler endpoints.SpotifyHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audio-feature-processing",
		Short: "Gets recent track information by combining information from LastFM and Spotify.",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				fileName = "audioFeaturesTimeSeries_20211014115400.txt"
			)

			fmt.Println("Starting the audio feature processing...")
			audioFeatures := parseInformationFromFile(fileName)
			getRecommendedTracks(audioFeatures, spotifyHandler)
		},
	}

	return cmd
}

func parseInformationFromFile(fileName string) []*models.TrackAudioFeatures {
	return models.GetTrackAudioFeatures(fileName)
}

func getRecommendedTracks(features []*models.TrackAudioFeatures, spotifyHandler endpoints.SpotifyHandler) {
	// Get stats for weekdays between 9-6
	work := getValuesBetweenTimeWeekdays(features, 9, 18)
	// Get recommended tracks for mondayWork
	recs := spotifyHandler.GetTrackRecommendations(work, "hardcore", 30)
	recIds := make([]spotify.ID, 0)
	for _, track := range recs.Tracks {
		recIds = append(recIds, track.ID)
	}
	spotifyHandler.CreatePlaylistAndAddTracks("Work Hardcore Mix", "", recIds)
}

func getValuesBetweenTimeWeekdays(features []*models.TrackAudioFeatures, startHour, endHour int) *models.AttributeStats {
	acousticness := make([]float64, 0)
	danceability := make([]float64, 0)
	energy := make([]float64, 0)
	instrumentalness := make([]float64, 0)
	liveness := make([]float64, 0)
	loudness := make([]float64, 0)
	speechiness := make([]float64, 0)
	tempo := make([]float64, 0)
	valence := make([]float64, 0)

	for _, t := range features {
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
