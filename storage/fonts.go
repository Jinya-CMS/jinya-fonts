package storage

import (
	"fmt"
	"github.com/gosimple/slug"
	"jinya-fonts/database"
)

func GetFontFileName(name, weight, style, fileType string, googleFont bool) string {
	prefix := "google"
	if !googleFont {
		prefix = "custom"
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s", prefix, slug.Make(name), weight, style, fileType)
}

func GetFontFiles(name string) ([]database.File, error) {
	return database.Select[database.File]("select * from file where font = $1 order by weight, style, type", name)
}

func GetFontFile(path string) ([]byte, error) {
	context, cancelFunc := getContext()
	redisClient, err := getRedisClient()
	if err == nil {
		defer cancelFunc()
		cachedFile := redisClient.Get(context, path)
		if bytes, err := cachedFile.Bytes(); err == nil {
			return bytes, nil
		}
	}

	return nil, nil
}

func setFontFile(font *database.Webfont, weight, style, fileType string, data []byte) (*database.File, error) {
	fileName := GetFontFileName(font.Name, weight, style, fileType, font.GoogleFont)

	path := fmt.Sprintf("/fonts/%s", fileName)
	metadata := database.File{
		Path:   path,
		Weight: weight,
		Style:  style,
		Type:   fileType,
		Font:   font.Name,
	}

	_ = addCachedFontFile(font.Name, weight, style, fileType, data, font.GoogleFont)

	_, err := database.GetDbMap().Exec("insert into file (path, weight, style, type, font) values ($1, $2, $3, $4, $5) on conflict do nothing", path, weight, style, fileType, font.Name)

	return &metadata, err
}

func removeFontFile(font *database.Webfont, weight, style, fileType string) error {
	_ = removeCachedFontFile(font.Name, weight, style, fileType, font.GoogleFont)

	fileName := GetFontFileName(font.Name, weight, style, fileType, font.GoogleFont)

	_, err := database.GetDbMap().Exec("delete from file where path = $1", fileName)

	return err
}

func SetFontFile(name, weight, style, fileType string, data []byte) (*database.File, error) {
	font, err := database.Get[database.Webfont](name)
	if err != nil {
		return nil, err
	}

	if font.GoogleFont {
		return nil, fmt.Errorf("cannot set file on google font")
	}

	return setFontFile(font, weight, style, fileType, data)
}

func RemoveFontFile(name, weight, style, fileType string) error {
	font, err := database.Get[database.Webfont](name)
	if err != nil {
		return err
	}

	if font.GoogleFont {
		return fmt.Errorf("cannot remove file from google font")
	}

	return removeFontFile(font, weight, style, fileType)
}

func SetGoogleFontFile(name, weight, style, fileType string, data []byte) (*database.File, error) {
	font, err := database.Get[database.Webfont](name)
	if err != nil {
		return nil, err
	}

	return setFontFile(font, weight, style, fileType, data)
}

func RemoveGoogleFontFile(name, weight, style, fileType string) error {
	font, err := database.Get[database.Webfont](name)
	if err != nil {
		return err
	}

	return removeFontFile(font, weight, style, fileType)
}
