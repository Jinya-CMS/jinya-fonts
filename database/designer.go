package database

import (
	"fmt"
)

type Designer struct {
	Name string `yaml:"name" json:"name"`
	Bio  string `yaml:"bio" json:"bio"`
}

func GetDesigners(name string) ([]Designer, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	return font.Designers, nil
}

func CreateDesigner(name string, designerName string, bio string) (*Designer, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	designer := Designer{
		Name: designerName,
		Bio:  bio,
	}

	font.Designers = append(font.Designers, designer)
	err = writeFont(*font)

	return &designer, err
}

func DeleteDesigner(name string, designerName string) error {
	font, err := GetFont(name)
	if err != nil {
		return err
	}

	if !font.GoogleFont {
		return fmt.Errorf("cannot delete google font")
	}

	var designers []Designer
	for _, item := range font.Designers {
		if item.Name != name {
			designers = append(designers, item)
		}
	}

	font.Designers = designers

	return writeFont(*font)
}
