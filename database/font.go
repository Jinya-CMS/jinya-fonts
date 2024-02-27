package database

import (
	"gopkg.in/yaml.v3"
	"jinya-fonts/config"
	"os"
)

type Metadata struct {
	Path         string `yaml:"path" json:"path"`
	Subset       string `yaml:"subset" json:"subset"`
	UnicodeRange string `yaml:"unicode_range" json:"unicodeRange"`
	Weight       string `yaml:"weight" json:"weight"`
	Style        string `yaml:"style" json:"style"`
	Category     string `yaml:"category" json:"category"`
	FontName     string `yaml:"-" json:"-"`
}

type Designer struct {
	Name string `yaml:"name" json:"name"`
	Bio  string `yaml:"bio" json:"bio"`
}

type Webfont struct {
	Name        string     `yaml:"name" json:"name"`
	Fonts       []Metadata `yaml:"fonts" json:"fonts"`
	Description string     `yaml:"description,omitempty" json:"description"`
	Designers   []Designer `yaml:"designers,omitempty" json:"designers"`
	License     string     `yaml:"license,omitempty" json:"license"`
	Category    string     `yaml:"category,omitempty" json:"category"`
	GoogleFont  bool       `yaml:"google_font" json:"-"`
}

func writeFont(webfont Webfont) error {
	yamlFileData, err := yaml.Marshal(&webfont)
	if err != nil {
		return err
	}

	return os.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+webfont.Name+".yaml", yamlFileData, 0644)
}

func GetAllFonts() ([]Webfont, error) {
	files, err := os.ReadDir(config.LoadedConfiguration.FontFileFolder)
	if err != nil {
		return nil, err
	}

	var availableFonts []Webfont

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		yamlFileData, err := os.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + file.Name())
		if err != nil {
			continue
		}

		fontFile := new(Webfont)
		err = yaml.Unmarshal(yamlFileData, fontFile)
		if err != nil {
			continue
		}

		availableFonts = append(availableFonts, *fontFile)
	}

	return availableFonts, err
}

func GetFont(name string) (*Webfont, error) {
	yamlFile := config.LoadedConfiguration.FontFileFolder + "/" + name + ".yaml"

	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		return nil, err
	}

	yamlFileData, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	fontFile := new(Webfont)
	err = yaml.Unmarshal(yamlFileData, fontFile)

	return fontFile, err
}

func CreateFont(name string, description string, license string, category string) (*Webfont, error) {
	webfont := Webfont{
		Name:        name,
		License:     license,
		Category:    category,
		Description: description,
		GoogleFont:  false,
	}

	if err := writeFont(webfont); err != nil {
		return nil, err
	}

	return &webfont, nil
}

func UpdateFont(name string, description string, license string, category string) error {
	webfont, err := GetFont(name)
	if err != nil {
		return err
	}

	webfont.Description = description
	webfont.License = license
	webfont.Category = category

	return writeFont(*webfont)
}

func DeleteFont(name string) error {
	return os.Remove(config.LoadedConfiguration.FontFileFolder + "/" + name + ".yaml")
}

func CreateGoogleFont(webfont *Webfont) error {
	return writeFont(*webfont)
}
