package tvdb

// Summary struct
type Summary struct {
	AiredEpisodes string   `json:"airedEpisodes"`
	AiredSeasons  []string `json:"airedSeasons"`
	DvdEpisodes   string   `json:"dvdEpisodes"`
	DvdSeasons    []string `json:"dvdSeasons"`
}
