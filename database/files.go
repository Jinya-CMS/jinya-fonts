package database

import (
	"fmt"
	"jinya-fonts/config"
	"os"
)

type Metadata struct {
	Path     string `yaml:"path" json:"path"`
	Weight   string `yaml:"weight" json:"weight"`
	Style    string `yaml:"style" json:"style"`
	Category string `yaml:"category" json:"category"`
}

func GetFontFiles(name string) ([]Metadata, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	return font.Fonts, nil
}

func AddFontFile(name string, data []byte, weight string, style string) (*Metadata, error) {
	font, err := GetFont(name)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+name, 0775)
	if err != nil {
		return nil, err
	}

	filename := name + "." + weight + "." + style + ".woff2"
	path := config.LoadedConfiguration.FontFileFolder + "/" + name + "/" + filename
	err = os.WriteFile(path, data, 0775)
	if err != nil {
		return nil, err
	}

	metadata := Metadata{
		Path:   filename,
		Weight: weight,
		Style:  style,
	}

	font.Fonts = append(font.Fonts, metadata)
	err = writeFont(*font)

	return &metadata, err
}

func UpdateFontFile(name string, data []byte, weight string, style string) error {
	font, err := GetFont(name)
	if err != nil {
		return err
	}

	exists := false
	for _, metadata := range font.Fonts {
		if metadata.Style == style && metadata.Weight == weight {
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("not found")
	}

	err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+name, 0775)
	if err != nil {
		return err
	}

	filename := config.LoadedConfiguration.FontFileFolder + "/" + name + "/" + name + "." + weight + "." + style + "." + ".woff2"
	err = os.WriteFile(filename, data, 0664)
	if err != nil {
		return err
	}

	return err
}

func DeleteFontFile(name string, weight string, style string) error {
	font, err := GetFont(name)
	if err != nil {
		return err
	}

	if !font.GoogleFont {
		return fmt.Errorf("cannot delete google font")
	}

	err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+name, 0775)
	if err != nil {
		return err
	}

	filename := config.LoadedConfiguration.FontFileFolder + "/" + name + "/" + name + "." + weight + "." + style + "." + ".woff2"
	err = os.Remove(filename)
	if err != nil {
		return err
	}

	var fonts []Metadata
	for _, item := range font.Fonts {
		if item.Path != filename {
			fonts = append(fonts, item)
		}
	}

	font.Fonts = fonts

	return writeFont(*font)
}
