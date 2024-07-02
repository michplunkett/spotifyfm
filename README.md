# SpotifyFM
This application is a `Golang` CLI that is in a bit of a transient state. Since I am currently in school, it has been put on the back-burner and is currently something I contribute to when the time feels right.

## Project Requirements
- Golang v1.22
- A [Spotify API token](https://developer.spotify.com/documentation/web-api) and the following environment variables:
  - `SPOTIFY_ID`: The application ID created from the above web-API link.
  - `SPOTIFY_SECRET`: The application secret created from the above web-API link.
  - `SPOTIFY_USER_NAME`: Your username on Spotify.
- A [Last.FM token](https://www.last.fm/api/authentication) and the following environment variables:
  - `LAST_FM_API_KEY`: The application ID created from the above web-API link.
  - `LAST_FM_SHARED_SECRET`: The application secret created from the above web-API link.

## Commands
- `go run main.go lastfm`: Holds commands that use the Spotify endpoints.
- `go run main.go spotify`: Holds commands that use the Last.FM endpoints.
- `go run main.go both`: Holds commands that use both the Spotify and Last.FM endpoints.
