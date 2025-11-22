package services

import (
	"encoding/json"
	"io"
	"mediaserver/views"
	"net/http"
	"os"
)

var OmdbApiUrl = "http://www.omdbapi.com/"
var OmdbApiKey = os.Getenv("OMDB_API_KEY")

func GetMoviesBySearch(title string) ([]views.OmdbSearchMovie, error) {
	resp, err := http.Get(OmdbApiUrl + "?s=" + title + "&apikey=" + OmdbApiKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var movieData views.OmdbSearchResult
	if err := json.Unmarshal(body, &movieData); err != nil {
		return nil, err
	}
	movies := movieData.Search
	err = json.Unmarshal(body, &movieData)
	if err != nil {
		return nil, err
	}

	return movies, nil

}
