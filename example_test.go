package tvdb_test

import (
	"fmt"
	"net/url"
	"os"
	"sort"

	"github.com/pioz/tvdb"
)

func ExampleClient_Login() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
}

func ExampleClient_SearchByName() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	res, err := c.SearchByName("Game of Thrones")
	if err != nil {
		panic(err)
	}
	fmt.Println(res[0].SeriesName)
	// Output: Game of Thrones
}

func ExampleClient_BestSearch() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		// The request response is a 404: this means no results have been found
		if tvdb.HaveCodeError(404, err) {
			fmt.Println("Series not found")
		} else {
			panic(err)
		}
	}
	fmt.Println(series.SeriesName)
	// Output: Game of Thrones
}

func ExampleClient_GetSeriesEpisodes() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		panic(err)
	}
	err = c.GetSeriesEpisodes(&series, url.Values{"airedSeason": {"2"}}) // params can be nil to retrieve all episodes
	if err != nil {
		panic(err)
	}
	fmt.Println(series.GetEpisode(2, 10).EpisodeName)
	// Output: Valar Morghulis
}

func ExampleClient_GetSeriesFanartImages() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		panic(err)
	}
	err = c.GetSeriesFanartImages(&series)
	if err != nil {
		panic(err)
	}
	url := tvdb.ImageURL(series.Images[0].FileName)
	fmt.Println(url)
	// Output: https://thetvdb.com/banners/fanart/original/121361-3.jpg
}

func ExampleSeries_GetSeasonEpisodes() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		panic(err)
	}
	err = c.GetSeriesEpisodes(&series, nil)
	if err != nil {
		panic(err)
	}
	episodes := series.GetSeasonEpisodes(1)
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].AiredEpisodeNumber < episodes[j].AiredEpisodeNumber
	})
	for _, ep := range episodes {
		fmt.Printf("1x%02d: %s - ", ep.AiredEpisodeNumber, ep.EpisodeName)
	}
	// Output: 1x01: Winter Is Coming - 1x02: The Kingsroad - 1x03: Lord Snow - 1x04: Cripples, Bastards, and Broken Things - 1x05: The Wolf and the Lion - 1x06: A Golden Crown - 1x07: You Win or You Die - 1x08: The Pointy End - 1x09: Baelor - 1x10: Fire and Blood -
}

func ExampleSeries_GetEpisode() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Username: os.Getenv("TVDB_USERNAME"), Userkey: os.Getenv("TVDB_USERKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		panic(err)
	}
	err = c.GetSeriesEpisodes(&series, nil)
	if err != nil {
		panic(err)
	}
	// Print the title of the episode 4x08 (season 4, episode 8)
	fmt.Println(series.GetEpisode(4, 8).EpisodeName)
	// Output: The Mountain and the Viper
}

func ExampleHaveCodeError() {
	c := tvdb.Client{Apikey: "WRONG APIKEY"}
	err := c.Login()
	if err == nil {
		panic("Impossible!")
	}
	fmt.Println(tvdb.HaveCodeError(401, err))
	// Output: true
}
