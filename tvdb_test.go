package tvdb

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const APIKEY string = "6F6E61197C18C895"

func TestLogin(t *testing.T) {
	c := login(t)
	assert.NotEmpty(t, c.token)
}

func TestRefreshToken(t *testing.T) {
	c := login(t)
	err := c.RefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, c.token)
}

func TestSearch(t *testing.T) {
	c := login(t)
	res, err := c.SearchByName("Game of Thrones")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "Game of Thrones", res[0].SeriesName)
}

func TestSearchByImdbID(t *testing.T) {
	c := login(t)
	res, err := c.SearchByImdbID("tt0944947")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "Game of Thrones", res[0].SeriesName)
}

func TestBestSearch(t *testing.T) {
	c := login(t)
	res, err := c.BestSearch("Game of Thrones")
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, Error404(err))
	assert.Equal(t, "Game of Thrones", res.SeriesName)
	res, err = c.BestSearch("kajdsfhasdkjhfsadkjhfasdkh")
	assert.True(t, Error404(err))
}

func TestGetSeries(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeries(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tt0944947", s.ImdbID)
}

func TestGetSeriesActors(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesActors(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(s.Actors))
	assert.Equal(t, "Ian McElhinney", s.Actors[0].Name)
}

func TestGetSeriesEpisodes(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesEpisodes(&s, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(s.Episodes))
	assert.Equal(t, 8, len(s.Episodes))
	assert.Equal(t, 10, len(s.Episodes[1]))
	assert.Equal(t, 10, len(s.Episodes[2]))
	assert.Equal(t, 10, len(s.Episodes[3]))
	assert.Equal(t, 10, len(s.Episodes[4]))
	assert.Equal(t, 10, len(s.Episodes[5]))
	assert.Equal(t, 10, len(s.Episodes[6]))
	assert.Equal(t, 7, len(s.Episodes[7]))
	assert.Equal(t, "Winter Is Coming", s.Episodes[1][1].EpisodeName)
	assert.Equal(t, "The Mountain and the Viper", s.Episodes[4][8].EpisodeName)
	assert.Equal(t, "The Dragon and the Wolf", s.Episodes[7][7].EpisodeName)
	err = c.GetSeriesEpisodes(&s, url.Values{"airedSeason": {"2"}})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(s.Episodes))
	assert.Nil(t, s.Episodes[1])
}

func TestGetSeriesSummary(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesSummary(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "105", s.Summary.AiredEpisodes)
	assert.Equal(t, 8, len(s.Summary.AiredSeasons))
}

func TestGetSeriesPosterImages(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesPosterImages(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "posters/121361-1.jpg", s.Images[0].FileName)
	assert.Equal(t, ImageURL(s.Images[0].FileName), "https://thetvdb.com/banners/posters/121361-1.jpg")
}

func login(t *testing.T) Client {
	c := Client{Apikey: APIKEY, Language: "en"}
	err := c.Login()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func getSerie(t *testing.T, c Client, name string) Series {
	series, err := c.BestSearch(name)
	if err != nil {
		t.Fatal(err)
	}
	return series
}
