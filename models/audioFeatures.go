package models

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TrackAudioFeatures struct {
	ListenDate       time.Time
	Name             string
	Album            string
	Artist           string
	Duration         float64
	SpotifyID        string
	Acousticness     float64
	Danceability     float64
	Energy           float64
	Instrumentalness float64
	Liveness         float64
	Loudness         float64
	Speechiness      float64
	Tempo            float64
	Valence          float64
}

func GetTrackAudioFeatures(fileName string) []*TrackAudioFeatures {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("failed to open audio features file")
	}
	audioFeatures := make([]*TrackAudioFeatures, 0)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	skipTheFirstLine := false
	for scanner.Scan() {
		// Checkin' a bool is quicker than string checkin'
		if !skipTheFirstLine && strings.HasPrefix(scanner.Text(), "ListenDate") {
			skipTheFirstLine = true
			continue
		}

		features := audioFeatureTextLineToDomain(scanner.Text())
		if features == nil {
			log.Println("could not process this line: ", scanner.Text())
			continue
		}

		audioFeatures = append(audioFeatures, features)
	}

	return audioFeatures
}

func audioFeatureTextLineToDomain(featureString string) *TrackAudioFeatures {
	f := strings.Split(featureString, "\t")
	if len(f) != 15 {
		fmt.Println("the length is ", len(f))
		return nil
	}

	listenDate, _ := time.Parse(time.RFC3339, f[0])
	duration, _ := strconv.ParseFloat(f[4], 64)
	acousticness, _ := strconv.ParseFloat(f[6], 64)
	danceability, _ := strconv.ParseFloat(f[7], 64)
	energy, _ := strconv.ParseFloat(f[8], 64)
	instrumentalness, _ := strconv.ParseFloat(f[9], 64)
	liveness, _ := strconv.ParseFloat(f[10], 64)
	loudness, _ := strconv.ParseFloat(f[11], 64)
	speechiness, _ := strconv.ParseFloat(f[12], 64)
	tempo, _ := strconv.ParseFloat(f[13], 64)
	valence, _ := strconv.ParseFloat(f[14], 64)
	return &TrackAudioFeatures{
		ListenDate:       listenDate,
		Name:             f[1],
		Album:            f[2],
		Artist:           f[3],
		Duration:         duration,
		SpotifyID:        f[5],
		Acousticness:     acousticness,
		Danceability:     danceability,
		Energy:           energy,
		Instrumentalness: instrumentalness,
		Liveness:         liveness,
		Loudness:         loudness,
		Speechiness:      speechiness,
		Tempo:            tempo,
		Valence:          valence,
	}
}

type BasicStats struct {
	Max    float64
	Median float64
	Min    float64
}

type AttributeStats struct {
	Acousticness     *BasicStats
	Danceability     *BasicStats
	Energy           *BasicStats
	Instrumentalness *BasicStats
	Liveness         *BasicStats
	Loudness         *BasicStats
	Speechiness      *BasicStats
	Tempo            *BasicStats
	Valence          *BasicStats
}

func BuildBasicStats(attr []float64) *BasicStats {
	sort.SliceStable(attr, func(i, j int) bool { return attr[i] < attr[j] })
	sliceLen := len(attr)

	stats := &BasicStats{
		Max:    attr[sliceLen-1],
		Median: 0,
		Min:    attr[0],
	}

	if sliceLen%2 == 0 {
		stats.Median = (attr[sliceLen/2] + attr[sliceLen/2+1]) / 2.00
	} else {
		stats.Median = attr[sliceLen/2+1]
	}

	return stats
}
