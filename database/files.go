package database

import (
	"cmp"
	"fmt"
	"github.com/gosimple/slug"
	"slices"
)

type File struct {
	Path   string `json:"path" bson:"path"`
	Weight string `json:"weight" bson:"weight"`
	Style  string `json:"style" bson:"style"`
	Data   []byte `json:"-" bson:"data"`
	Type   string `json:"type" bson:"type"`
}

func GetFontFileName(name, weight, style, fileType string, googleFont bool) string {
	prefix := "google"
	if !googleFont {
		prefix = "custom"
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s", prefix, slug.Make(name), weight, style, fileType)
}

func GetFontFiles(name string) ([]File, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	var fonts []File
	for _, file := range font.Fonts {
		fonts = append(fonts, file)
	}

	slices.SortFunc(fonts, func(a, b File) int {
		return cmp.Compare(a.Weight+"."+a.Style, b.Weight+"."+b.Style)
	})

	return fonts, nil
}

func SetFontFile(name, weight, style, fileType string, data []byte) (*File, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	if font.GoogleFont {
		return nil, fmt.Errorf("cannot set file on google font")
	}

	fileName := GetFontFileName(name, weight, style, fileType, font.GoogleFont)
	path := fmt.Sprintf("/fonts/%s", fileName)
	metadata := File{
		Path:   path,
		Weight: weight,
		Style:  style,
		Type:   fileType,
		Data:   data,
	}

	_ = AddCachedFontFile(name, weight, style, fileType, data, false)

	if font.Fonts == nil {
		font.Fonts = map[string]File{}
	}

	font.Fonts[fileName] = metadata
	err = UpdateFont(font)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func RemoveFontFile(name, weight, style, fileType string) error {
	font, err := GetFont(name)
	if err != nil {
		return err
	}

	if font.GoogleFont {
		return fmt.Errorf("cannot remove file from google font")
	}

	_ = RemoveCachedFontFile(name, weight, style, fileType, false)

	delete(font.Fonts, GetFontFileName(name, weight, style, fileType, font.GoogleFont))

	return UpdateFont(font)
}
