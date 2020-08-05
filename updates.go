package tvdb

//Update Specifies when a show was last updated
type Update struct {
	//ID Show Identifier
	ID int `json:"id"`
	//LastUpdated epoch date when the show was updated
	LastUpdated int `json:"lastUpdated"`
}
