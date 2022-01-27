package meta

import (
	"gopkg.in/yaml.v3"
	"jinya-fonts/config"
	"log"
	"os"
	"sync"
)

type FontFileMeta struct {
	Path         string `yaml:"path" json:"path"`
	Subset       string `yaml:"subset" json:"subset"`
	UnicodeRange string `yaml:"unicode_range" json:"unicodeRange"`
	Weight       string `yaml:"weight" json:"weight"`
	Style        string `yaml:"style" json:"style"`
	FontName     string `yaml:"-" json:"-"`
}

type FontDesigner struct {
	Name string `yaml:"name" json:"name"`
	Bio  string `yaml:"bio" json:"bio"`
}

type FontFile struct {
	Name        string         `yaml:"name" json:"name"`
	Fonts       []FontFileMeta `yaml:"fonts" json:"fonts"`
	Description string         `yaml:"description,omitempty" json:"description"`
	Designers   []FontDesigner `yaml:"designers,omitempty" json:"designers"`
	License     string         `yaml:"license,omitempty" json:"license"`
	Category    string         `yaml:"category,omitempty" json:"category"`
	GoogleFont  bool           `yaml:"google_font" json:"-"`
}

var fontWriteMetadataMutex = sync.Mutex{}

func SaveFontFileMetadata(file FontFile) error {
	configuration := config.LoadedConfiguration
	data, err := yaml.Marshal(file)
	if err != nil {
		log.Printf("Failed to marshal font meta data %s", file.Name)
		return err
	}

	log.Printf("Lock fontWriteMetadataMutex for font %s", file.Name)
	fontWriteMetadataMutex.Lock()
	err = os.WriteFile(configuration.FontFileFolder+"/"+file.Name+".yaml", data, 0775)
	if err != nil {
		log.Printf("Failed to save font meta data %s", file.Name)
	}
	log.Printf("Unlock fontWriteMetadataMutex for font %s", file.Name)
	fontWriteMetadataMutex.Unlock()

	return err
}

func LoadFontFileCache(name string) (*FontFile, error) {
	configuration := config.LoadedConfiguration
	fontPath := configuration.FontFileFolder + "/" + name + ".yaml"
	file, err := os.Open(fontPath)
	if err != nil {
		return nil, err
	}

	fontFile := new(FontFile)
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(fontFile)
	if err != nil {
		return nil, err
	}
	return fontFile, nil
}
