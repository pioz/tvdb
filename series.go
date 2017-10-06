package tvdb

// Series struct store all data of an episode.
type Series struct {
	Added           string   `json:"added"`
	AddedBy         int      `json:"addedBy"`
	AirsDayOfWeek   string   `json:"airsDayOfWeek"`
	AirsTime        string   `json:"airsTime"`
	Aliases         []string `json:"aliases"`
	Banner          string   `json:"banner"`
	FirstAired      string   `json:"firstAired"`
	Genre           []string `json:"genre"`
	ID              int      `json:"id"`
	ImdbID          string   `json:"imdbId"`
	LastUpdated     int      `json:"lastUpdated"`
	Network         string   `json:"network"`
	NetworkID       string   `json:"networkId"`
	Overview        string   `json:"overview"`
	Rating          string   `json:"rating"`
	Runtime         string   `json:"runtime"`
	SeriesID        string   `json:"seriesId"`
	SeriesName      string   `json:"seriesName"`
	SiteRating      float32  `json:"siteRating"`
	SiteRatingCount int      `json:"siteRatingCount"`
	Status          string   `json:"status"`
	Zap2itID        string   `json:"zap2itId"`
	// Slice of the series actors, filled with GetSeriesActors method.
	Actors []Actor
	// Slice of the series episodes, filled with GetSeriesEpisodes method.
	Episodes []Episode
	// Slice of the series summary, filled with GetSeriesSummary method.
	Summary Summary
	// Slice of the series images.
	Images []Image
}

// Empty verify if the series's fields are empty and don't are filled by an api
// response.
func (s *Series) Empty() bool {
	return s.ID == 0 && s.SeriesName == ""
}

// BannerURL returns the image banner url of the series.
func (s *Series) BannerURL() string {
	return ImageURL(s.Banner)
}

// GetSeasonEpisodes select and returns the episodes of the series by season
// number.
func (s *Series) GetSeasonEpisodes(season int) []*Episode {
	episodes := make([]*Episode, 0)
	for i := range s.Episodes {
		if s.Episodes[i].AiredSeason == season {
			episodes = append(episodes, &s.Episodes[i])
		}
	}
	return episodes
}

// GetEpisode select and returns a specific episode of the series by season and
// episode number.
func (s *Series) GetEpisode(season, number int) *Episode {
	for i := range s.Episodes {
		if s.Episodes[i].AiredSeason == season && s.Episodes[i].AiredEpisodeNumber == number {
			return &s.Episodes[i]
		}
	}
	return nil
}
