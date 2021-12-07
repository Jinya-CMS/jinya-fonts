package meta

import (
	"gopkg.in/yaml.v3"
	"jinya-fonts/config"
	"log"
	"os"
	"sync"
)

type FontFileMeta struct {
	Path         string `yaml:"path"`
	Subset       string `yaml:"subset"`
	Variant      string `yaml:"variant"`
	UnicodeRange string `yaml:"unicode_range"`
	Weight       string `yaml:"weight"`
	Style        string `yaml:"style"`
	Category     string `yaml:"category"`
	FontName     string `yaml:"-"`
}

type FontDesigner struct {
	Name string `yaml:"name"`
	Bio  string `yaml:"bio"`
}

type FontFile struct {
	Name        string         `yaml:"name"`
	Fonts       []FontFileMeta `yaml:"fonts"`
	Description string         `yaml:"description,omitempty"`
	Designers   []FontDesigner `yaml:"designers,omitempty"`
	License     string         `yaml:"license,omitempty"`
	Category    string         `yaml:"category,omitempty"`
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
