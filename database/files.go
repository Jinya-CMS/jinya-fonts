package database

import (
	"bytes"
	"cmp"
	"fmt"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"slices"
)

type File struct {
	Path   string             `json:"path" bson:"path"`
	Weight string             `json:"weight" bson:"weight"`
	Style  string             `json:"style" bson:"style"`
	Type   string             `json:"type" bson:"type"`
	FileId primitive.ObjectID `json:"-" bson:"fileId"`
}

func GetFontFileName(name, weight, style, fileType string, googleFont bool) string {
	prefix := "google"
	if !googleFont {
		prefix = "custom"
	}

	return fmt.Sprintf("%s.%s.%s.%s.%s", prefix, slug.Make(name), weight, style, fileType)
}

func GetFontFileData(file File) (io.Reader, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	bucket, err := getFontFileBucket(client)
	if err != nil {
		return nil, err
	}

	return bucket.OpenDownloadStream(file.FileId)
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

func SetFontFile(name, weight, style, fileType string, data []byte, force bool) (*File, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	if font.GoogleFont && !force {
		return nil, fmt.Errorf("cannot set file on google font")
	}

	fileName := GetFontFileName(name, weight, style, fileType, font.GoogleFont)

	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	bucket, err := getFontFileBucket(client)
	if err != nil {
		return nil, err
	}

	stream, err := bucket.UploadFromStream(fileName, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/fonts/%s", fileName)
	metadata := File{
		Path:   path,
		Weight: weight,
		Style:  style,
		Type:   fileType,
		FileId: stream,
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

	client, err := openConnection()
	if err != nil {
		return err
	}

	defer closeConnection(client)

	bucket, err := getFontFileBucket(client)
	if err != nil {
		return err
	}

	fileName := GetFontFileName(name, weight, style, fileType, font.GoogleFont)
	if file, exists := font.Fonts[fileName]; exists {
		err = bucket.Delete(file.FileId)
		if err != nil {
			return err
		}
	}

	delete(font.Fonts, fileName)

	return UpdateFont(font)
}
