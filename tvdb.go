// Package tvdb is cool
package tvdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client struct
type Client struct {
	Apikey   string
	Language string
	token    string
	client   http.Client
}

// BaseURL is the endpoint URL
const BaseURL string = "https://api.thetvdb.com"

// Login perform a post request to retrieve the token from the apikey
func (c *Client) Login() error {
	resp, err := c.performPOSTRequest("/login", map[string]string{"apikey": c.Apikey})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(loginAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	c.token = data.Token
	return nil
}

// RefreshToken refresh the token
func (c *Client) RefreshToken() error {
	resp, err := c.performGETRequest("/refresh_token", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(loginAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	c.token = data.Token
	return nil
}

// GetLanguages returns all avaiable languages
func (c *Client) GetLanguages() ([]Language, error) {
	resp, err := c.performGETRequest("/languages", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data := new(languagesAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

// SearchByName allows the user to search for a series based on the serie name.
func (c *Client) SearchByName(q string) ([]Series, error) {
	return c.search(url.Values{"name": {q}})
}

// SearchByImdbID allows the user to search for a series based on the serie name.
func (c *Client) SearchByImdbID(q string) ([]Series, error) {
	return c.search(url.Values{"imdbId": {q}})
}

// SearchByZap2itID allows the user to search for a series based on the serie name.
func (c *Client) SearchByZap2itID(q string) ([]Series, error) {
	return c.search(url.Values{"zap2itId": {q}})
}

// BestSearch return the best result
func (c *Client) BestSearch(q string) (Series, error) {
	res, err := c.SearchByName(q)
	if err != nil {
		return Series{}, err
	}
	for _, series := range res {
		if strings.ToLower(series.SeriesName) == strings.ToLower(q) {
			return series, nil
		}
	}
	for _, serie := range res {
		for _, alias := range serie.Aliases {
			if strings.ToLower(alias) == strings.ToLower(q) {
				return serie, nil
			}
		}
	}
	return res[0], nil
}

// GetSeries retrieve all series's fields
func (c *Client) GetSeries(s *Series) error {
	if s.Empty() {
		return errors.New("The serie is empty")
	}
	resp, err := c.performGETRequest(fmt.Sprintf("/series/%d", s.ID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(seriesAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	*s = data.Data
	return nil
}

// GetSeriesActors retrieve all series's actors
func (c *Client) GetSeriesActors(s *Series) error {
	if s.Empty() {
		return errors.New("The serie is empty")
	}
	resp, err := c.performGETRequest(fmt.Sprintf("/series/%d/actors", s.ID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(actorsAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	s.Actors = data.Data
	return nil
}

// GetSeriesEpisodes retrieve all series's actors
func (c *Client) GetSeriesEpisodes(s *Series, params url.Values) error {
	if s.Empty() {
		return errors.New("The serie is empty")
	}
	episodes := make([]Episode, 0)
	var (
		data episodesAPIResponse
		resp *http.Response
		err  error
	)
	if params == nil {
		params = url.Values{"page": {"1"}}
	}
	for page := 1; ; page++ {
		params.Set("page", strconv.Itoa(page))
		resp, err = c.performGETRequest(fmt.Sprintf("/series/%d/episodes/query", s.ID), params)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		err = parseResponse(resp.Body, &data)
		if err != nil {
			return err
		}
		episodes = append(episodes, data.Data...)
		if len(data.Data) < 100 {
			break
		}
	}
	s.Episodes = arrangeEpisodes(episodes)
	return nil
}

// GetSeriesSummary retrieve all series's actors
func (c *Client) GetSeriesSummary(s *Series) error {
	if s.Empty() {
		return errors.New("The serie is empty")
	}
	resp, err := c.performGETRequest(fmt.Sprintf("/series/%d/episodes/summary", s.ID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(summaryAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	s.Summary = data.Data
	return nil
}

// GetEpisode retrieve all episode's fields
func (c *Client) GetEpisode(e *Episode) error {
	if e.Empty() {
		return errors.New("The episode is empty")
	}
	resp, err := c.performGETRequest(fmt.Sprintf("/episodes/%d", e.ID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(episodeAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	*e = data.Data
	return nil
}

// GetSeriesFanartImages return all fanart images of the series
func (c *Client) GetSeriesFanartImages(s *Series) error {
	return c.getSeriesImages(s, "fanart")
}

// GetSeriesPosterImages return all poster images of the series
func (c *Client) GetSeriesPosterImages(s *Series) error {
	return c.getSeriesImages(s, "poster")
}

// GetSeriesSeasonImages return all season images of the series
func (c *Client) GetSeriesSeasonImages(s *Series) error {
	return c.getSeriesImages(s, "season")
}

// GetSeriesSeasonwideImages return all seasonwide images of the series
func (c *Client) GetSeriesSeasonwideImages(s *Series) error {
	return c.getSeriesImages(s, "seasonwide")
}

// GetSeriesSeriesImages return all series images of the series
func (c *Client) GetSeriesSeriesImages(s *Series) error {
	return c.getSeriesImages(s, "series")
}

func (c *Client) getSeriesImages(s *Series, keyType string) error {
	if s.Empty() {
		return errors.New("The serie is empty")
	}
	resp, err := c.performGETRequest(fmt.Sprintf("/series/%d/images/query", s.ID), url.Values{"keyType": {keyType}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data := new(imagesAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return err
	}
	s.Images = data.Data
	return nil
}

func (c *Client) search(params url.Values) ([]Series, error) {
	resp, err := c.performGETRequest("/search/series", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data := new(searchAPIResponse)
	err = parseResponse(resp.Body, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (c *Client) performGETRequest(path string, params url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", BaseURL, path), nil)
	req.URL.RawQuery = params.Encode()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", c.Language)
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	resp, err := c.client.Do(req)
	if err == nil && resp.StatusCode != 200 {
		return nil, &RequestError{resp.StatusCode}
	}
	return resp, err
}

func (c *Client) performPOSTRequest(path string, params map[string]string) (*http.Response, error) {
	json, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", BaseURL, path), bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", c.Language)
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	resp, err := c.client.Do(req)
	if err == nil && resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Get a response with status code %d", resp.StatusCode)
		return nil, errors.New(errorMessage)
	}
	return resp, err
}

func parseResponse(body io.ReadCloser, data interface{}) error {
	// b, err := ioutil.ReadAll(body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(b))
	// err = json.Unmarshal(b, &data)
	err := json.NewDecoder(body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func arrangeEpisodes(episodes []Episode) Episodes {
	m := make(Episodes)
	for _, ep := range episodes {
		if m[ep.AiredSeason] == nil {
			m[ep.AiredSeason] = map[int]Episode{ep.AiredEpisodeNumber: ep}
		} else {
			m[ep.AiredSeason][ep.AiredEpisodeNumber] = ep
		}
	}
	return m
}
