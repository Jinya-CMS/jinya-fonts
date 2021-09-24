package meta

import (
	"gopkg.in/yaml.v3"
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
	Type         string `yaml:"type"`
	FontName     string `yaml:"-"`
}

type FontFile struct {
	Name  string          `yaml:"name"`
	Fonts []*FontFileMeta `yaml:"fonts"`
}

var fontWriteMetadataMutex = sync.Mutex{}

func SaveFontFileMetadata(name string, dataDir string, metaData []*FontFileMeta) error {
	file := FontFile{
		Name:  name,
		Fonts: metaData,
	}
	data, err := yaml.Marshal(file)
	if err != nil {
		log.Printf("Failed to marshal font meta data %s", name)
		return err
	}

	log.Printf("Lock fontWriteMetadataMutex for font %s", name)
	fontWriteMetadataMutex.Lock()
	err = os.WriteFile(dataDir+"/"+name+".yaml", data, 0775)
	if err != nil {
		log.Printf("Failed to save font meta data %s", name)
	}
	log.Printf("Unlock fontWriteMetadataMutex for font %s", name)
	fontWriteMetadataMutex.Unlock()

	return err
}
