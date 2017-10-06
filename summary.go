package tvdb

// Summary struct store all data of a summary.
type Summary struct {
	AiredEpisodes string   `json:"airedEpisodes"`
	AiredSeasons  []string `json:"airedSeasons"`
	DvdEpisodes   string   `json:"dvdEpisodes"`
	DvdSeasons    []string `json:"dvdSeasons"`
}
