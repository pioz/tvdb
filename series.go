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
	Episodes        Episodes
	Summary         Summary
	Images          []Image
}

// Empty return true if the Serie's fields are empty
func (s *Series) Empty() bool {
	return s.ID == 0 && s.SeriesName == ""
}
