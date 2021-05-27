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
	Duration         float32
	SpotifyID        string
	Acousticness     float32
	Danceability     float32
	Energy           float32
	Instrumentalness float32
	Liveness         float32
	Loudness         float32
	Speechiness      float32
	Tempo            float32
	Valence          float32
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
	duration, _ := strconv.ParseFloat(f[4], 32)
	acousticness, _ := strconv.ParseFloat(f[6], 32)
	danceability, _ := strconv.ParseFloat(f[7], 32)
	energy, _ := strconv.ParseFloat(f[8], 32)
	instrumentalness, _ := strconv.ParseFloat(f[9], 32)
	liveness, _ := strconv.ParseFloat(f[10], 32)
	loudness, _ := strconv.ParseFloat(f[11], 32)
	speechiness, _ := strconv.ParseFloat(f[12], 32)
	tempo, _ := strconv.ParseFloat(f[13], 32)
	valence, _ := strconv.ParseFloat(f[14], 32)
	return &TrackAudioFeatures{
		ListenDate:       listenDate,
		Name:             f[1],
		Album:            f[2],
		Artist:           f[3],
		Duration:         float32(duration),
		SpotifyID:        f[5],
		Acousticness:     float32(acousticness),
		Danceability:     float32(danceability),
		Energy:           float32(energy),
		Instrumentalness: float32(instrumentalness),
		Liveness:         float32(liveness),
		Loudness:         float32(loudness),
		Speechiness:      float32(speechiness),
		Tempo:            float32(tempo),
		Valence:          float32(valence),
	}
}

type BasicStats struct {
	Max    float32
	Median float32
	Min    float32
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

func BuildBasicStats(attr []float32) *BasicStats {
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
