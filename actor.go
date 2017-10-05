package tvdb

// Actor struct
type Actor struct {
	ID          int    `json:"id"`
	Image       string `json:"image"`
	ImageAdded  string `json:"imageAdded"`
	ImageAuthor int    `json:"imageAuthor"`
	LastUpdated string `json:"lastUpdated"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	SeriesID    int    `json:"seriesId"`
	SortOrder   int    `json:"sortOrder"`
}
