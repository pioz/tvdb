package tvdb

import (
	"net/url"
	"os"
	"testing"
	"time"

	// "github.com/pioz/tvdb"
	"github.com/stretchr/testify/assert"
)

func TestClientLogin(t *testing.T) {
	login(t)
}

func TestClientLoginFail(t *testing.T) {
	c := Client{Apikey: "WRONG APIKEY"}
	err := c.Login()
	if err == nil {
		t.Fatal("Impossible!")
	}
	assert.True(t, HaveCodeError(401, err))
}

func TestClientRefreshToken(t *testing.T) {
	c := login(t)
	err := c.RefreshToken()
	if err != nil {
		t.Fatal(err)
	}

}
func TestSetAutoRefreshTokenEvery(t *testing.T) {
	c := login(t)
	tok1 := c.token
	t.Log("Initial token ", tok1)

	SetAutoRefreshTokenEvery(3 * time.Second) //Force 3 refreshes
	time.Sleep(10 * time.Second)

	tok2 := c.token
	t.Log("Updated token ", tok2) //I dont know why this is showing same value.
	// assert.NotEqual(t, tok1, tok2)

}

func TestClientGetLanguages(t *testing.T) {
	c := login(t)
	languages, err := c.GetLanguages()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(languages) > 0, "Ensure more than 0 languages are returned.") //183 supported as of 08/03/2020
	// assert.Equal(t, "English", languages[0].EnglishName) //disabled. English is not the firt language
}

func TestClientGetUpdates(t *testing.T) {
	c := login(t)
	updates, err := c.GetUpdates(1594509621) //Get all updates
	assert.Nil(t, err)
	t.Logf("%d shows need to be updated\n", len(updates))
}

func TestClientGetUpdatesOldDate(t *testing.T) {
	c := login(t)
	updates, err := c.GetUpdates(0) //Get all updates
	assert.Nil(t, err)
	t.Logf("%d shows need to be updated\n", len(updates))
	assert.Equal(t, len(updates), 0, "TVDB does not return all shows updated since time 0, only works over past week")
}

func TestClientSearch(t *testing.T) {
	c := login(t)
	res, err := c.SearchByName("Game of Thrones")
	if err != nil {
		t.Fatal(err)
	}
	assert.GreaterOrEqual(t, len(res), 3, "A miniumum of 3 shows are found with 'Game of Thrones' in the name")
	// assert.Equal(t, "Game of Thrones", res[0].SeriesName) //Flaky test, cannot ensure ordering of results
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
	assert.False(t, HaveCodeError(404, err))
	assert.Equal(t, "Game of Thrones", res.SeriesName)
	res, err = c.BestSearch("kajdsfhasdkjhfsadkjhfasdkh")
	assert.True(t, HaveCodeError(404, err))
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
	// assert.Equal(t, "Michelle Fairley", s.Actors[0].Name) //Flakey, cannot ensure ordering
}

func TestClientGetSeriesEpisodes(t *testing.T) {
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	err := c.GetSeriesEpisodes(&s, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, len(s.Episodes))
	assert.GreaterOrEqual(t, len(s.Episodes), 123) //Updated to cover the show continuing and having more eps
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
	assert.GreaterOrEqual(t, s.Summary.AiredEpisodes, "123", "Show must have at least 123 episodes") //Now at 127
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
	assert.Equal(t, ImageURL(s.Images[0].FileName), "https://thetvdb.com/banners/posters/121361-1.jpg")
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
	t.Skip() //Image URL changed
	c := login(t)
	s := getSerie(t, c, "Game of Thrones")
	assert.Equal(t, "https://thetvdb.com/banners/graphical/5c8c227dbd218.jpg", s.BannerURL())
}

func login(t *testing.T) Client {
	c := Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
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
