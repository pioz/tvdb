// Package tvdb is a wrapper in Go for the TVDB json api version 2. With this
// package you can make http requests to https://api.thetvdb.com.
//
// You can install this package with:
//
//   $ go get github.com/pioz/tvdb
//
// See also
//
// https://api.thetvdb.com/swagger for TVDB api version 2 documentation.
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

// Client does the work of perform the REST requests to the TVDB api endpoints.
// With its methods you can run almost all the requests provided in the TVDB
// api.
type Client struct {
	// The TVDB API key, User key, User name. You can get them here http://thetvdb.com/?tab=apiregister
	Apikey string
	Userkey string
	Username string
	// The language with which you want to obtain the data (if not set english is
	// used)
	Language string
	token    string
	client   http.Client
}

// BaseURL where the TVDB api is accessible.
const BaseURL string = "https://api.thetvdb.com"

// Login is used to retrieve a valid token which will be used to make any other
// requests to the TVDB api. The token is stored in the Client struct.
func (c *Client) Login() error {
	resp, err := c.performPOSTRequest("/login", map[string]string{"apikey": c.Apikey, "userkey": c.Userkey, "username": c.Username})
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

// RefreshToken is used to refresh the current token.
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

// GetLanguages returns all avaiable languages, a slice of Language.
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

// SearchByName allows to search for a series based on the series name. Returns
// the slice of the series found.
func (c *Client) SearchByName(q string) ([]Series, error) {
	return c.search(url.Values{"name": {q}})
}

// SearchByImdbID allows to search for a series based on the IMDB id
// (https://www.imdb.com). Returns the slice of the series found.
func (c *Client) SearchByImdbID(q string) ([]Series, error) {
	return c.search(url.Values{"imdbId": {q}})
}

// SearchByZap2itID allows to search for a series based on the Zap2it id
// (http://zap2it.com). Returns the slice of the series found.
func (c *Client) SearchByZap2itID(q string) ([]Series, error) {
	return c.search(url.Values{"zap2itId": {q}})
}

// BestSearch returns the best Series based on the name (q).
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

// GetSeries retrieve all series's fields. If a series is returned from a search
// method it will not have all fields filled. This method fills all fields of
// the series passed by reference as parameter.
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

// GetSeriesActors retrieve all series's actors. Actors slice is accessible from
// series.Actors struct field.
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

// GetSeriesEpisodes retrieve series's episodes. Episodes slice is accessible
// from series.Episodes struct field but is better obtain episodes using the
// series's methods GetEpisodes and GetEpisode. The parameter params is the
// parameters used to perform the request to the api. Valid params are:
// absoluteNumber, airedSeason, airedEpisode, dvdSeason, dvdEpisode, imdbId,
// page (100 episodes per page, if page is not passed retrieve all episodes).
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
	s.Episodes = episodes
	return nil
}

// GetSeriesSummary retrieve the summary of the episodes and seasons available
// for the series. Summary is accessible from series.Summary struct field.
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

// GetEpisode retrieve all episode's fields. If an episode is returned from the
// GetEpisodes method it will not have all fields filled. This method fills all
// fields of the episode passed by reference as parameter.
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

// GetSeriesFanartImages retrieve fanart images of the series. These images are
// accessible from series.Images struct field.
func (c *Client) GetSeriesFanartImages(s *Series) error {
	return c.getSeriesImages(s, "fanart")
}

// GetSeriesPosterImages retrieve poster images of the series. These images are
// accessible from series.Images struct field.
func (c *Client) GetSeriesPosterImages(s *Series) error {
	return c.getSeriesImages(s, "poster")
}

// GetSeriesSeasonImages retrieve season images of the series. These images are
// accessible from series.Images struct field.
func (c *Client) GetSeriesSeasonImages(s *Series) error {
	return c.getSeriesImages(s, "season")
}

// GetSeriesSeasonwideImages retrieve season wide images of the series. These images are
// accessible from series.Images struct field.
func (c *Client) GetSeriesSeasonwideImages(s *Series) error {
	return c.getSeriesImages(s, "seasonwide")
}

// GetSeriesSeriesImages retrieve series images of the series. These images are
// accessible from series.Images struct field.
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
		return nil, &RequestError{resp.StatusCode}
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

// func arrangeEpisodes(episodes []Episode) Episodes {
// 	m := make(Episodes)
// 	for _, ep := range episodes {
// 		if m[ep.AiredSeason] == nil {
// 			m[ep.AiredSeason] = map[int]Episode{ep.AiredEpisodeNumber: ep}
// 		} else {
// 			m[ep.AiredSeason][ep.AiredEpisodeNumber] = ep
// 		}
// 	}
// 	return m
// }
