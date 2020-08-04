package tvdb

//Updates contains an slice of show ids with their last updated date in Epoch Time
type Updates struct {
	UpdateEntry []UpdateEntry `json:"data"`
}

//UpdateEntry Specifies when a show was last updated
type UpdateEntry struct {
	//ID Show Identifier
	ID int `json:"id"`
	//LastUpdated epoch date when the show was updated
	LastUpdated int `json:"lastUpdated"`
}
