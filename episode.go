package tvdb

// Episode struct
type Episode struct {
	AbsoluteNumber     int      `json:"absoluteNumber"`
	AiredEpisodeNumber int      `json:"airedEpisodeNumber"`
	AiredSeason        int      `json:"airedSeason"`
	AirsAfterSeason    int      `json:"airsAfterSeason"`
	AirsBeforeEpisode  int      `json:"airsBeforeEpisode"`
	AirsBeforeSeason   int      `json:"airsBeforeSeason"`
	Director           string   `json:"director"`
	Directors          []string `json:"directors"`
	DvdChapter         int      `json:"dvdChapter"`
	DvdDiscid          string   `json:"dvdDiscid"`
	DvdEpisodeNumber   int      `json:"dvdEpisodeNumber"`
	DvdSeason          int      `json:"dvdSeason"`
	EpisodeName        string   `json:"episodeName"`
	Filename           string   `json:"filename"`
	FirstAired         string   `json:"firstAired"`
	GuestStars         []string `json:"guestStars"`
	ID                 int      `json:"id"`
	ImdbID             string   `json:"imdbId"`
	LastUpdated        int      `json:"lastUpdated"`
	LastUpdatedBy      string   `json:"lastUpdatedBy"`
	Overview           string   `json:"overview"`
	ProductionCode     string   `json:"productionCode"`
	SeriesID           string   `json:"seriesID"`
	ShowURL            string   `json:"showURL"`
	SiteRating         int      `json:"siteRating"`
	SiteRatingCount    int      `json:"siteRatingCount"`
	ThumbAdded         string   `json:"thumbAdded"`
	ThumbAuthor        int      `json:"thumbAuthor"`
	ThumbHeight        string   `json:"thumbHeight"`
	ThumbWidth         string   `json:"thumbWidth"`
	Writers            []string `json:"writers"`
}

// Episodes rappresent a map of episodes organized by season and episode number
type Episodes map[int]map[int]Episode
