package fontsync

import (
	"encoding/json"
	"fmt"
	"io"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"jinya-fonts/storage"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

type GoogleWebfont struct {
	Family       string            `json:"family"`
	Variants     []string          `json:"variants"`
	Subsets      []string          `json:"subsets"`
	Version      string            `json:"version"`
	LastModified string            `json:"lastModified"`
	Files        map[string]string `json:"files"`
	Category     string            `json:"category"`
	Kind         string            `json:"kind"`
	Menu         string            `json:"menu"`
}

type SyncJob struct{}

func (s SyncJob) Run() {
	err := Sync()
	if err != nil {
		log.Printf("Failed to run sync job %s", err.Error())
	}
}

func downloadWoff2FontList() ([]GoogleWebfont, error) {
	log.Println("Download font list")
	settings, err := database.GetSettings()
	if err != nil {
		return nil, err
	}

	familyFilter := strings.Join(settings.FilterByName, "&family=")
	if len(familyFilter) > 0 {
		familyFilter = "&family=" + strings.ReplaceAll(familyFilter, " ", "+")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://webfonts.googleapis.com/v1/webfonts?capability=WOFF2&key=%s%s", config.LoadedConfiguration.ApiKey, familyFilter), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch list")
	}

	log.Println("Got result from webfonts.googleapis.com")

	var webFontList struct {
		Items []GoogleWebfont `json:"items"`
	}

	decoder := json.NewDecoder(res.Body)

	log.Println("Decode result")
	err = decoder.Decode(&webFontList)
	if err != nil {
		return nil, err
	}

	return webFontList.Items, nil
}

func downloadFontFiles(ttf bool, font GoogleWebfont, cpu int) {
	fileType := "woff2"
	if ttf == true {
		fileType = "ttf"
	}

	for weightAndStyle, file := range font.Files {
		weight := "400"
		style := "normal"
		if strings.HasSuffix(weightAndStyle, "italic") {
			style = "italic"
		}
		if weightAndStyle != "italic" && weightAndStyle != "regular" {
			weight = weightAndStyle[0:3]
		}

		response, err := http.Get(file)
		if err != nil {
			logWithCpu(cpu, "Failed to load font file from Google server: %s", err.Error())
			continue
		}

		if response.StatusCode != http.StatusOK {
			logWithCpu(cpu, "Failed to load font file from Google server: %s", response.Status)
			continue
		}

		fontData, err := io.ReadAll(response.Body)
		if err != nil {
			logWithCpu(cpu, "Failed to read font file data: %s", err.Error())
			continue
		}

		go func() {
			_, err := storage.SetGoogleFontFile(font.Family, weight, style, fileType, fontData)
			if err != nil {
				logWithCpu(cpu, "Failed to save font file: %s", err.Error())
			}
		}()
	}
}

func handleWebfont(channel chan GoogleWebfont, wg *sync.WaitGroup, cpu int) {
	for font := range channel {
		webfontMetadata, err := getGoogleWebfontMetadata(cpu, font.Family)
		if err != nil {
			logWithCpu(cpu, "Failed to load webfont metadata: %s", err.Error())
			continue
		}

		fontFile := database.Webfont{
			Name:        font.Family,
			Description: webfontMetadata.Description,
			License:     webfontMetadata.License,
			Category:    webfontMetadata.Category,
			GoogleFont:  true,
		}

		err = database.Insert(fontFile)
		if err != nil {
			logWithCpu(cpu, "Failed to create new font: %s", err.Error())
			continue
		}

		for _, designer := range webfontMetadata.Designers {
			designer.Font = font.Family
			err = database.Insert(designer)
			if err != nil {
				logWithCpu(cpu, "Failed to create new font: %s", err.Error())
			}
		}

		downloadFontFiles(false, font, cpu)
	}

	wg.Done()
}

func Sync() error {
	log.Println("Grab font list")
	fonts, err := downloadWoff2FontList()
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	fontChannel := make(chan GoogleWebfont)
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleWebfont(fontChannel, wg, i)
	}

	err = database.ClearGoogleFonts()
	if err != nil {
		log.Printf("Failed to clear google fonts: %s", err.Error())
	}

	storage.ClearGoogleFontsCache()

	for _, font := range fonts {
		fontChannel <- font
	}

	close(fontChannel)
	wg.Wait()

	return nil
}
