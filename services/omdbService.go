package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mediaserver/data"
	"mediaserver/views"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	OmdbApiUrl = "http://www.omdbapi.com/"
	OmdbApiKey = os.Getenv("OMDB_API_KEY")

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)

func GetMoviesBySearch(title string) ([]views.OmdbSearchMovie, error) {
	if OmdbApiKey == "" {
		return nil, errors.New("OMDB_API_KEY is not set")
	}

	escaped := url.QueryEscape(title)
	reqUrl := fmt.Sprintf("%s?s=%s&apikey=%s", OmdbApiUrl, escaped, OmdbApiKey)

	resp, err := httpClient.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("omdb returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var movieData views.OmdbSearchResult
	if err := json.Unmarshal(body, &movieData); err != nil {
		return nil, err
	}

	return movieData.Search, nil
}

func GetOmdbMovieById(id string) (data.OmdbMovie, error) {
	if OmdbApiKey == "" {
		return data.OmdbMovie{}, errors.New("OMDB_API_KEY is not set")
	}

	escaped := url.QueryEscape(id)
	reqUrl := fmt.Sprintf("%s?i=%s&apikey=%s", OmdbApiUrl, escaped, OmdbApiKey)

	resp, err := httpClient.Get(reqUrl)
	if err != nil {
		return data.OmdbMovie{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data.OmdbMovie{}, fmt.Errorf("omdb returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data.OmdbMovie{}, err
	}

	var movie data.OmdbMovie
	if err := json.Unmarshal(body, &movie); err != nil {
		return data.OmdbMovie{}, err
	}

	return movie, nil
}

func GetMovieViewById(id string) (*views.OmdbMovie, error) {
	movie, err := GetOmdbMovieById(id)
	if err != nil {
		return nil, err
	}

	return movie.ToView(), nil
}
