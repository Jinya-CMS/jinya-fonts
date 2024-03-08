package fontsync

import (
	"encoding/json"
	"fmt"
	"io"
	"jinya-fonts/config"
	"jinya-fonts/database"
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

func downloadFontFiles(ttf bool, font GoogleWebfont, cpu int) (fontFiles map[string]database.File, fontData map[string][]byte, err error) {
	fontFiles = map[string]database.File{}
	fontData = map[string][]byte{}

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

		fileName := database.GetFontFileName(font.Family, weight, style, fileType, true)
		fontData[fileName], err = io.ReadAll(response.Body)
		if err != nil {
			logWithCpu(cpu, "Failed to read font file data: %s", err.Error())
			continue
		}

		path := fmt.Sprintf("/fonts/%s", fileName)

		fontFile := database.File{
			Path:   path,
			Weight: weight,
			Style:  style,
			Type:   fileType,
		}

		if err != nil {
			logWithCpu(cpu, "Failed to add font file: %s", err.Error())
			continue
		}

		go database.AddCachedFontFile(font.Family, weight, style, fileType, fontData[fileName], true)

		fontFiles[fileName] = fontFile
	}

	return
}

func handleWebfont(channel chan GoogleWebfont, wg *sync.WaitGroup, cpu int) {
	for font := range channel {
		webfontMetadata, err := getGoogleWebfontMetadata(cpu, font.Family)
		if err != nil {
			logWithCpu(cpu, "Failed to load webfont metadata: %s", err.Error())
			continue
		}

		fonts, files, err := downloadFontFiles(false, font, cpu)
		if err != nil {
			logWithCpu(cpu, "Failed to create new font: %s", err.Error())
			continue
		}

		fontFile := database.Webfont{
			Name:        font.Family,
			Description: webfontMetadata.Description,
			Designers:   webfontMetadata.Designers,
			License:     webfontMetadata.License,
			Category:    webfontMetadata.Category,
			GoogleFont:  true,
			Fonts:       fonts,
		}

		err = database.CreateFont(&fontFile)
		if err != nil {
			logWithCpu(cpu, "Failed to create new font: %s", err.Error())
			continue
		}

		for p, file := range files {
			meta := fonts[p]
			_, err := database.SetFontFile(font.Family, meta.Weight, meta.Style, meta.Type, file, true)
			if err != nil {
				logWithCpu(cpu, "Failed to save font file: %s", err.Error())
				continue
			}
		}
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

	database.ClearGoogleFonts()
	database.ClearGoogleFontsCache()

	for _, font := range fonts {
		fontChannel <- font
	}

	close(fontChannel)
	wg.Wait()

	return nil
}
