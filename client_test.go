package tvdb_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/pioz/tvdb"
	"github.com/stretchr/testify/assert"
)

func TestClientLogin(t *testing.T) {
	login(t)
}

func TestClientLoginFail(t *testing.T) {
	c := tvdb.Client{Apikey: "WRONG APIKEY"}
	err := c.Login()
	if err == nil {
		t.Fatal("Impossible!")
	}
	assert.True(t, tvdb.HaveCodeError(401, err))
}

func TestClientRefreshToken(t *testing.T) {
	c := login(t)
	err := c.RefreshToken()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientGetLanguages(t *testing.T) {
	c := login(t)
	languages, err := c.GetLanguages()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 23, len(languages))
	assert.Equal(t, "English", languages[0].EnglishName)
}

func TestClientSearch(t *testing.T) {
	c := login(t)
	res, err := c.SearchByName("Game of Thrones")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "Game of Thrones", res[0].SeriesName)
}

func TestClientSearchByImdbID(t *testing.T) {
	c := login(t)
	res, err := c.SearchByImdbID("tt0944947")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "Game of Thrones", res[0].SeriesName)
}

func TestClientBestSearch(t *testing.T) {
	c := login(t)
	res, err := c.BestSearch("Game of Thrones")
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, tvdb.HaveCodeError(404, err))
	assert.Equal(t, "Game of Thrones", res.SeriesName)
	res, err = c.BestSearch("kajdsfhasdkjhfsadkjhfasdkh")
	assert.True(t, tvdb.HaveCodeError(404, err))
}

func TestClientGetSeries(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeries(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tt0944947", s.ImdbID)
}

func TestClientGetSeriesActors(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesActors(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(s.Actors))
	assert.Equal(t, "Michelle Fairley", s.Actors[0].Name)
}

func TestClientGetSeriesEpisodes(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesEpisodes(&s, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(s.Episodes))
	assert.Equal(t, 123, len(s.Episodes))
	assert.Equal(t, "Winter Is Coming", s.GetEpisode(1, 1).EpisodeName)
	assert.Equal(t, "The Mountain and the Viper", s.GetEpisode(4, 8).EpisodeName)
	assert.Equal(t, "The Dragon and the Wolf", s.GetEpisode(7, 7).EpisodeName)
	err = c.GetSeriesEpisodes(&s, url.Values{"airedSeason": {"2"}})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, len(s.Episodes))
	assert.Nil(t, s.GetEpisode(1, 1))
}

func TestClientGetSeriesSummary(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesSummary(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", s.Summary.AiredEpisodes)
	assert.Equal(t, 9, len(s.Summary.AiredSeasons))
}

func TestClientGetSeriesPosterImages(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesPosterImages(&s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "posters/121361-1.jpg", s.Images[0].FileName)
	assert.Equal(t, tvdb.ImageURL(s.Images[0].FileName), "https://thetvdb.com/banners/posters/121361-1.jpg")
}

func TestClientGetEpisode(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesEpisodes(&s, url.Values{"airedSeason": {"1"}})
	if err != nil {
		t.Fatal(err)
	}
	err = c.GetEpisode(s.GetEpisode(1, 1))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "tt1480055", s.GetEpisode(1, 1).ImdbID)
}

func TestSeriesGetSeasonEpisodes(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	assert.Equal(t, 0, len(s.GetSeasonEpisodes(2)))
	err := c.GetSeriesEpisodes(&s, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, len(s.GetSeasonEpisodes(2)))
}

func TestSeriesBannerURL(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	assert.Equal(t, "https://thetvdb.com/banners/graphical/5c8c227dbd218.jpg", s.BannerURL())
}

func login(t *testing.T) tvdb.Client {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func getSerie(t *testing.T, c tvdb.Client, name string) tvdb.Series {
	series, err := c.BestSearch(name)
	if err != nil {
		t.Fatal(err)
	}
	return series
}
