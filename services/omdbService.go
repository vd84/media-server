package services

import (
	"encoding/json"
	"io"
	"mediaserver/data"
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

func GetMovieById(id string) (data.OmdbMovie, error) {
	resp, err := http.Get(OmdbApiUrl + "?i=" + id + "&apikey=" + OmdbApiKey)
	if err != nil {
		return data.OmdbMovie{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data.OmdbMovie{}, err
	}

	var movie data.OmdbMovie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		return data.OmdbMovie{}, err
	}

	return movie, nil
}

func GetMovieViewById(id string) (*views.OmdbMovie, error) {
	movie, err := GetMovieById(id)
	if err != nil {
		return nil, err
	}

	return movie.ToView(), nil
}
