package tvdb

// Image struct store all data of an image.
type Image struct {
	FileName    string `json:"fileName"`
	ID          int    `json:"id"`
	KeyType     string `json:"keyType"`
	LanguageID  int    `json:"languageId"`
	Resolution  string `json:"resolution"`
	SubKey      string `json:"subKey"`
	Thumbnail   string `json:"thumbnail"`
	RatingsInfo Rating `json:"ratingsInfo"`
}

// holds image ratings
type Rating struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}

// ImageURL returns the complete URL of an image. This because the images
// fileName returned by the TVDB api are relative. So this function simply join
// the base URL (https://thetvdb.com/banners) with the relative path passed as
// parameter.
func ImageURL(fileName string) string {
	return "https://thetvdb.com/banners/" + fileName
}
