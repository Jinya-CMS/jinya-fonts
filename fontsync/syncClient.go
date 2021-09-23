package fontsync

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"jinya-fonts/config"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
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
}

type WebFont struct {
	Family       string            `json:"family"`
	Variants     []string          `json:"variants"`
	Subsets      []string          `json:"subsets"`
	Version      string            `json:"version"`
	LastModified string            `json:"lastModified"`
	Files        map[string]string `json:"files"`
	Category     string            `json:"category"`
	Kind         string            `json:"kind"`
}

type WebFontList struct {
	Items []WebFont `json:"items"`
	Kind  string    `json:"kind"`
}

const (
	FontTypeTtf   = "Mozilla/5.0"
	FontTypeWoff  = "Mozilla/4.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1667.0 Safari/537.36"
	FontTypeWoff2 = "Mozilla/5.0 (Windows NT 6.3; rv:39.0) Gecko/20100101 Firefox/44.0"
)

var (
	fontFolderCreationMutex = sync.Mutex{}
	fontWriteMetadataMutex  = sync.Mutex{}
	copyFontDataMutex       = sync.Mutex{}
)

func downloadFontList(apiKey string) ([]WebFont, error) {
	log.Println("Download font list")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://webfonts.googleapis.com/v1/webfonts?key=%s", apiKey), nil)
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

	var webFontList WebFontList

	decoder := json.NewDecoder(res.Body)

	log.Println("Decode result")
	err = decoder.Decode(&webFontList)
	if err != nil {
		return nil, err
	}

	return webFontList.Items, nil
}

type fontDownloadJob struct {
	Name      string
	Subset    string
	Variant   string
	UserAgent string
	Type      string
}

func saveFontFile(configuration *config.Configuration, channel chan fontDownloadJob, wg *sync.WaitGroup, idx int) {
	for job := range channel {
		log.Printf("CPU %d: Download font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://fonts.googleapis.com/css?subset=%s&family=%s:%s", job.Subset, job.Name, job.Variant), nil)
		if err != nil {
			log.Printf("CPU %d: Failed to create request for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			continue
		}

		req.Header.Add("User-Agent", job.UserAgent)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("CPU %d: Failed to get data for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			continue
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("CPU %d: Failed to get data for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			continue
		}

		fontCss, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("CPU %d: Failed to read body for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			continue
		}

		splitFaces := strings.Split(string(fontCss), "}")

		for _, face := range splitFaces {
			log.Printf("CPU %d: Find font face url %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			fontFaceRegex := regexp.MustCompile(`src: url\((?P<font>.*)\) `)
			fontFaceMatches := fontFaceRegex.FindStringSubmatch(face)
			if len(fontFaceMatches) != 2 {
				log.Printf("CPU %d: Failed to find url for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			fontIndex := fontFaceRegex.SubexpIndex("font")
			fontUrl := fontFaceMatches[fontIndex]

			log.Printf("CPU %d: Find font unicode range %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			unicodeRangeRegex := regexp.MustCompile(`unicode-range: (?P<range>.*);`)
			rangeMatches := unicodeRangeRegex.FindStringSubmatch(face)
			rangeIndex := unicodeRangeRegex.SubexpIndex("range")
			rangeValue := ""
			if rangeIndex != -1 && job.Type == "woff2" {
				rangeValue = rangeMatches[rangeIndex]
			}

			weightRegex := regexp.MustCompile(`font-weight: (?P<weight>.*);`)

			log.Printf("CPU %d: Find font weight %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			weightMatches := weightRegex.FindStringSubmatch(face)
			if len(weightMatches) != 2 {
				log.Printf("CPU %d: Failed to find font-weight for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			weightIndex := weightRegex.SubexpIndex("weight")
			weightValue := weightMatches[weightIndex]

			log.Printf("CPU %d: Find font style %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			styleRegex := regexp.MustCompile(`font-style: (?P<style>.*);`)
			styleMatches := styleRegex.FindStringSubmatch(face)
			if len(styleMatches) != 2 {
				log.Printf("CPU %d: Failed to find font-style for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			styleIndex := styleRegex.SubexpIndex("style")
			styleValue := styleMatches[styleIndex]

			subsetValue := job.Subset
			if job.Type == "woff2" {
				log.Printf("CPU %d: Find font subset %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				subsetRegex := regexp.MustCompile(`\/\* (?P<subset>.*) \*\/`)
				subsetMatches := subsetRegex.FindStringSubmatch(face)
				if len(styleMatches) != 2 {
					log.Printf("CPU %d: Failed to find font-subset for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
					continue
				}

				subsetIndex := subsetRegex.SubexpIndex("subset")
				subsetValue = subsetMatches[subsetIndex]
			}

			res, err = http.Get(fontUrl)
			if err != nil {
				log.Printf("CPU %d: Failed to download font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			fontDir := configuration.FontFileFolder + "/" + job.Name
			log.Printf("CPU %d: Lock fontFolderCreationMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			fontFolderCreationMutex.Lock()
			err = os.MkdirAll(fontDir, 0775)
			if err != nil {
				log.Printf("CPU %d: Failed to download font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				fontFolderCreationMutex.Unlock()
				continue
			}
			log.Printf("CPU %d: Unlock fontFolderCreationMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			fontFolderCreationMutex.Unlock()

			log.Printf("CPU %d: Write font file %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			file, err := os.OpenFile(fontDir+"/"+job.Name+"."+subsetValue+"."+job.Variant+"."+job.Type, os.O_CREATE|os.O_WRONLY, 0775)
			if err != nil {
				log.Printf("CPU %d: Failed to open file to safe font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			log.Printf("CPU %d: Lock copyFontDataMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			copyFontDataMutex.Lock()
			_, err = io.Copy(file, res.Body)
			if err != nil {
				log.Printf("CPU %d: Failed to copy font into file %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				copyFontDataMutex.Lock()
				continue
			}
			log.Printf("CPU %d: Unlock copyFontDataMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			copyFontDataMutex.Unlock()

			meta := FontFileMeta{
				Path:         job.Name + "." + job.Type,
				Subset:       subsetValue,
				Variant:      job.Variant,
				UnicodeRange: rangeValue,
				Weight:       weightValue,
				Style:        styleValue,
				Type:         job.Type,
			}

			data, err := yaml.Marshal(meta)
			if err != nil {
				log.Printf("CPU %d: Failed to marshal font meta data %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				continue
			}

			log.Printf("CPU %d: Lock fontWriteMetadataMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			fontWriteMetadataMutex.Lock()
			err = os.WriteFile(fontDir+"/"+job.Name+"."+subsetValue+"."+job.Variant+"."+job.Type+".yaml", data, 0775)
			if err != nil {
				log.Printf("CPU %d: Failed to save font meta data %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
				fontWriteMetadataMutex.Unlock()
				continue
			}
			log.Printf("CPU %d: Unlock fontWriteMetadataMutex for font %s.%s %s %s", idx, job.Name, job.Type, job.Variant, job.Subset)
			fontWriteMetadataMutex.Unlock()
		}
	}

	wg.Done()
}

func Sync(configuration *config.Configuration) error {
	log.Println("Grab font list")
	fonts, err := downloadFontList(configuration.ApiKey)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	fontChannel := make(chan fontDownloadJob)
	for i := 0; i < runtime.NumCPU(); i++ {
		go saveFontFile(configuration, fontChannel, wg, i)
	}

	for _, font := range fonts {
		subsets := font.Subsets
		variants := font.Variants
		name := font.Family
		for _, subset := range subsets {
			for _, variant := range variants {
				fontChannel <- fontDownloadJob{
					Name:      name,
					Subset:    subset,
					Variant:   variant,
					UserAgent: FontTypeTtf,
					Type:      "ttf",
				}
				fontChannel <- fontDownloadJob{
					Name:      name,
					Subset:    subset,
					Variant:   variant,
					UserAgent: FontTypeWoff,
					Type:      "woff",
				}
				fontChannel <- fontDownloadJob{
					Name:      name,
					Subset:    subset,
					Variant:   variant,
					UserAgent: FontTypeWoff2,
					Type:      "woff2",
				}
			}
		}
	}

	close(fontChannel)
	wg.Wait()

	return nil
}
