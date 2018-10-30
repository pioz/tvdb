# TVDB API for Go

Golang wrapper for TVDB json api version 2.

The TVDB api version 2 documentation [can be found here](https://api.thetvdb.com/swagger).

## Installation

Install it yourself as:

    $ go get github.com/pioz/tvdb

(optional) To run unit tests:

    $ cd $GOPATH/src/github.com/pioz/tvdb
    $ TVDB_APIKEY=your_apikey TVDB_USERKEY=your_userkey TVDB_USERNAME=your_username go test -v

## Usage

First of all you need to get your API key, User key and User name:

* Register an account on http://thetvdb.com/?tab=register
* When you are logged register an api key on http://thetvdb.com/?tab=apiregister
* View your api key, user key and user name on http://thetvdb.com/?tab=userinfo

```Go
package main

import (
  "fmt"
  "github.com/pioz/tvdb"
)

func main() {
  c := tvdb.Client{Apikey: "YOUR API KEY", Userkey: "YOUR USER KEY", Username: "YOUR USER NAME"}
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
```

The complete __documentation__ can be found [here](https://godoc.org/github.com/pioz/tvdb).

## Missing REST endpoints

This wrapper do not coverage all 100% api REST endpoints.
Missing methods are:

* __Series__
    * filter: `GET /series/{id}/filter`
* __Updates__
    * updadad: `GET /updated/query`
* __Users__
    * user: `GET /user`
    * favorites: `GET /user/favorites`
    * delete favorites: `DELETE /user/favorites/{id}`
    * add favorites: `PUT /user/favorites/{id}`
    * ratings: `GET /user/ratings`
    * ratings with query: `GET /user/ratings/query`
    * delete rating: `DELETE /user/ratings/{itemType}/{itemId}`
    * add rating: `PUT /user/ratings/{itemType}/{itemId}/{itemRating}`

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/pioz/tvdb.

## License

The package is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
