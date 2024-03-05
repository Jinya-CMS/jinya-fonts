package database

import (
	"fmt"
	"slices"
)

type Designer struct {
	Name string `json:"name" bson:"name"`
	Bio  string `json:"bio" bson:"bio"`
}

func GetDesigners(name string) ([]Designer, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	return font.Designers, nil
}

func CreateDesigner(name string, designer Designer) (*Designer, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	if font.GoogleFont {
		return nil, fmt.Errorf("cannot add designers from a google font")
	}

	font.Designers = append(font.Designers, designer)
	err = UpdateFont(font)
	if err != nil {
		return nil, err
	}

	return &designer, err
}

func DeleteDesigner(name string, designerName string) error {
	font, err := GetFont(name)
	if err != nil {
		return err
	}

	if font.GoogleFont {
		return fmt.Errorf("cannot remove designers from a google font")
	}

	font.Designers = slices.DeleteFunc(font.Designers, func(designer Designer) bool {
		return designer.Name == designerName
	})

	return UpdateFont(font)
}
