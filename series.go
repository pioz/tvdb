package tvdb

// Series struct
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
	Actors          []Actor
	Episodes        []Episode
	Summary         Summary
	Images          []Image
}

// Empty return true if the Serie's fields are empty
func (s *Series) Empty() bool {
	return s.ID == 0 && s.SeriesName == ""
}

// BannerURL return the banner of the series
func (s *Series) BannerURL() string {
	return ImageURL(s.Banner)
}

// GetSeasonEpisodes return the episode of the series specific season
func (s *Series) GetSeasonEpisodes(season int) []*Episode {
	episodes := make([]*Episode, 0)
	for i := range s.Episodes {
		if s.Episodes[i].AiredSeason == season {
			episodes = append(episodes, &s.Episodes[i])
		}
	}
	return episodes
}

// GetEpisode return a specific episode of the series
func (s *Series) GetEpisode(season, number int) *Episode {
	for i := range s.Episodes {
		if s.Episodes[i].AiredSeason == season && s.Episodes[i].AiredEpisodeNumber == number {
			return &s.Episodes[i]
		}
	}
	return nil
}
