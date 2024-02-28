package fontsync

import (
	"encoding/json"
	"fmt"
	"io"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"log"
	"net/http"
	"os"
	"runtime"
	"slices"
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

func downloadFontList(apiKey string, fetchOnly []string) ([]GoogleWebfont, error) {
	log.Println("Download font list")
	familyFilter := strings.Join(fetchOnly, "&family=")
	if len(familyFilter) > 0 {
		familyFilter = "&family=" + strings.ReplaceAll(familyFilter, " ", "+")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://webfonts.googleapis.com/v1/webfonts?capability=WOFF2&key=%s%s", apiKey, familyFilter), nil)
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

func downloadFontFiles(font GoogleWebfont, cpu int) {
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

		data, err := io.ReadAll(response.Body)
		if err != nil {
			logWithCpu(cpu, "Failed to read font file data: %s", err.Error())
			continue
		}

		_, err = database.AddFontFile(font.Family, data, weight, style)
		if err != nil {
			logWithCpu(cpu, "Failed to add font file: %s", err.Error())
			continue
		}
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
			Designers:   webfontMetadata.Designers,
			License:     webfontMetadata.License,
			Category:    webfontMetadata.Category,
			GoogleFont:  true,
		}

		err = database.CreateGoogleFont(&fontFile)
		if err != nil {
			logWithCpu(cpu, "Failed to create new font: %s", err.Error())
			continue
		}

		downloadFontFiles(font, cpu)
	}

	wg.Done()
}

func Sync(configuration *config.Configuration) error {
	log.Printf("Create data directory if not existing: %s", configuration.FontFileFolder)
	err := os.MkdirAll(configuration.FontFileFolder, 0775)
	if err != nil {
		return err
	}

	log.Println("Grab font list")
	fonts, err := downloadFontList(configuration.ApiKey, configuration.FilterByName)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	fontChannel := make(chan GoogleWebfont)
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleWebfont(fontChannel, wg, i)
	}

	for _, font := range fonts {
		if len(configuration.FilterByName) > 0 && !slices.ContainsFunc(configuration.FilterByName, func(filter string) bool {
			return strings.ToLower(filter) == strings.ToLower(font.Family)
		}) {
			continue
		}

		fontChannel <- font
	}

	close(fontChannel)
	wg.Wait()

	return nil
}
