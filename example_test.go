package tvdb_test

import (
	"fmt"
	"net/url"
	"os"

	"github.com/pioz/tvdb"
)

func ExampleClient_Login() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
}

func ExampleClient_SearchByName() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
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
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
	err := c.Login()
	if err != nil {
		panic(err)
	}
	series, err := c.BestSearch("Game of Thrones")
	if err != nil {
		// The request response is a 404: this means no results have been found
		if tvdb.Error404(err) {
			fmt.Println("Series not found")
		} else {
			panic(err)
		}
	}
	fmt.Println(series.SeriesName)
	// Output: Game of Thrones
}

func ExampleClient_GetSeriesEpisodes() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
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
	fmt.Println(series.Episodes[2][10].EpisodeName)
	// Output: Valar Morghulis
}

func ExampleClient_GetSeriesFanartImages() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY"), Language: "en"}
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

func ExampleEpisodes() {
	c := tvdb.Client{Apikey: os.Getenv("TVDB_APIKEY")}
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
	fmt.Println(series.Episodes[4][8].EpisodeName)
	// Output: The Mountain and the Viper
}
