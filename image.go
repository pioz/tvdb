package tvdb

import (
	"fmt"
)

// Image struct
type Image struct {
	FileName   string `json:"fileName"`
	ID         int    `json:"id"`
	KeyType    string `json:"keyType"`
	LanguageID int    `json:"languageId"`
	Resolution string `json:"resolution"`
	SubKey     string `json:"subKey"`
	Thumbnail  string `json:"thumbnail"`
}

// ImageURL returns the complete URL of an image
func ImageURL(fileName string) string {
	return fmt.Sprintf("https://thetvdb.com/banners/%s", fileName)
}
