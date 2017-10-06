package tvdb

// Image struct store all data of an image.
type Image struct {
	FileName   string `json:"fileName"`
	ID         int    `json:"id"`
	KeyType    string `json:"keyType"`
	LanguageID int    `json:"languageId"`
	Resolution string `json:"resolution"`
	SubKey     string `json:"subKey"`
	Thumbnail  string `json:"thumbnail"`
}

// ImageURL returns the complete URL of an image. This because the images
// fileName returned by the TVDB api are relative. So this function simply join
// the base URL (https://thetvdb.com/banners) with the relative path passed as
// parameter.
func ImageURL(fileName string) string {
	return "https://thetvdb.com/banners/" + fileName
}
