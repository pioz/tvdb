package tvdb

// Language struct
type Language struct {
	Abbreviation string `json:"abbreviation"`
	EnglishName  string `json:"englishName"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}
