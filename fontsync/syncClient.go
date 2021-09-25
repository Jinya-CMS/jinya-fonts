package fontsync

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

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
	FontTypeWoff2 = "Mozilla/5.0 (Windows NT 6.3; rv:39.0) Gecko/20100101 Firefox/44.0"
)

var (
	fontFolderCreationMutex = sync.Mutex{}
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

type FontDownloadJob struct {
	Name     string
	Variant  string
	Category string
}

func saveFontFile(configuration *config.Configuration, channel chan []FontDownloadJob, wg *sync.WaitGroup, idx int) {
	for jobs := range channel {
		var fontData []meta.FontFileMeta
		var name string
		for _, job := range jobs {
			woff2Css, _ := fetchCss(idx, job, FontTypeWoff2)

			var faces []string
			faces = append(faces, strings.Split(string(woff2Css), "}")...)

			for _, face := range faces {
				if strings.Contains(face, "@font-face") {
					data, err := HandleFontFace(configuration, idx, job, face)
					data.Category = job.Category
					if err != nil {
						log.Printf("CPU %d: %s", idx, err.Error())
						continue
					}

					fontData = append(fontData, *data)
				}
			}
			name = job.Name
		}
		err := meta.SaveFontFileMetadata(name, configuration.FontFileFolder, fontData)
		if err != nil {
			log.Printf("CPU %d: %s", idx, err.Error())
		}
	}

	wg.Done()
}

func HandleFontFace(configuration *config.Configuration, idx int, job FontDownloadJob, face string) (*meta.FontFileMeta, error) {
	fontUrl, err := getFontFaceUrl(idx, job, face)
	if err != nil {
		return nil, err
	}

	rangeValue := getFontUnicodeRange(idx, job, face)
	weightValue, err := getFontWeight(idx, job, face)
	if err != nil {
		return nil, err
	}

	styleValue, err := getFontStyle(idx, job, face)
	if err != nil {
		return nil, err
	}

	subsetValue := "all"
	subsetValue, err = getFontSubset(idx, job, face, subsetValue)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(fontUrl)
	if err != nil {
		log.Printf("CPU %d: Failed to download font %s %s", idx, job.Name, job.Variant)
		return nil, err
	}

	fontDir := configuration.FontFileFolder + "/" + job.Name
	err = createFontDirectory(idx, job, err, fontDir)
	file, err := openFontFile(idx, job, err, fontDir, subsetValue)
	if err != nil {
		return nil, err
	}

	err = copyFontFileFromResponse(idx, job, file, res)
	if err != nil {
		return nil, err
	}

	path := job.Name + "." + subsetValue + "." + job.Variant + ".woff2"

	return &meta.FontFileMeta{
		Path:         path,
		Subset:       subsetValue,
		Variant:      job.Variant,
		UnicodeRange: rangeValue,
		Weight:       weightValue,
		Style:        styleValue,
		FontName:     job.Name,
	}, nil
}

func copyFontFileFromResponse(idx int, job FontDownloadJob, file *os.File, res *http.Response) error {
	log.Printf("CPU %d: Lock copyFontDataMutex for font %s %s", idx, job.Name, job.Variant)
	copyFontDataMutex.Lock()
	_, err := io.Copy(file, res.Body)
	if err != nil {
		log.Printf("CPU %d: Failed to copy font into file %s %s", idx, job.Name, job.Variant)
	}
	log.Printf("CPU %d: Unlock copyFontDataMutex for font %s %s", idx, job.Name, job.Variant)
	copyFontDataMutex.Unlock()
	return err
}

func openFontFile(idx int, job FontDownloadJob, err error, fontDir string, subsetValue string) (*os.File, error) {
	log.Printf("CPU %d: Open font file %s %s", idx, job.Name, job.Variant)
	file, err := os.OpenFile(fontDir+"/"+job.Name+"."+subsetValue+"."+job.Variant+".woff2", os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		log.Printf("CPU %d: Failed to open file to safe font %s %s", idx, job.Name, job.Variant)
		return nil, err
	}
	return file, nil
}

func createFontDirectory(idx int, job FontDownloadJob, err error, fontDir string) error {
	log.Printf("CPU %d: Lock fontFolderCreationMutex for font %s %s", idx, job.Name, job.Variant)
	fontFolderCreationMutex.Lock()
	err = os.MkdirAll(fontDir, 0775)
	if err != nil {
		log.Printf("CPU %d: Failed to download font %s %s", idx, job.Name, job.Variant)
	}
	log.Printf("CPU %d: Unlock fontFolderCreationMutex for font %s %s", idx, job.Name, job.Variant)
	fontFolderCreationMutex.Unlock()
	return err
}

func getFontSubset(idx int, job FontDownloadJob, face string, subsetValue string) (string, error) {
	log.Printf("CPU %d: Find font subset %s %s", idx, job.Name, job.Variant)
	subsetRegex := regexp.MustCompile(`\/\* (?P<subset>.*) \*\/`)
	subsetMatches := subsetRegex.FindStringSubmatch(face)
	if len(subsetMatches) != 2 {
		log.Printf("CPU %d: Failed to find font-subset for font %s %s", idx, job.Name, job.Variant)
		return "", fmt.Errorf("failed to find font subset")
	}

	subsetIndex := subsetRegex.SubexpIndex("subset")
	subsetValue = subsetMatches[subsetIndex]
	return subsetValue, nil
}

func getFontStyle(idx int, job FontDownloadJob, face string) (string, error) {
	log.Printf("CPU %d: Find font style %s %s", idx, job.Name, job.Variant)
	styleRegex := regexp.MustCompile(`font-style: (?P<style>.*);`)
	styleMatches := styleRegex.FindStringSubmatch(face)
	if len(styleMatches) != 2 {
		log.Printf("CPU %d: Failed to find font-style for font %s %s", idx, job.Name, job.Variant)
		return "", fmt.Errorf("failed to find font style")
	}

	styleIndex := styleRegex.SubexpIndex("style")
	styleValue := styleMatches[styleIndex]
	return styleValue, nil
}

func getFontWeight(idx int, job FontDownloadJob, face string) (string, error) {
	log.Printf("CPU %d: Find font weight %s %s", idx, job.Name, job.Variant)
	weightRegex := regexp.MustCompile(`font-weight: (?P<weight>.*);`)
	weightMatches := weightRegex.FindStringSubmatch(face)
	if len(weightMatches) != 2 {
		log.Printf("CPU %d: Failed to find font-weight for font %s %s", idx, job.Name, job.Variant)
		return "", fmt.Errorf("failed to find font weight")
	}

	weightIndex := weightRegex.SubexpIndex("weight")
	weightValue := weightMatches[weightIndex]
	return weightValue, nil
}

func getFontUnicodeRange(idx int, job FontDownloadJob, face string) string {
	log.Printf("CPU %d: Find font unicode range %s %s", idx, job.Name, job.Variant)
	unicodeRangeRegex := regexp.MustCompile(`unicode-range: (?P<range>.*);`)
	rangeMatches := unicodeRangeRegex.FindStringSubmatch(face)
	rangeIndex := unicodeRangeRegex.SubexpIndex("range")
	if rangeIndex != -1 {
		return rangeMatches[rangeIndex]
	}

	return ""
}

func getFontFaceUrl(idx int, job FontDownloadJob, face string) (string, error) {
	log.Printf("CPU %d: Find font face url %s %s", idx, job.Name, job.Variant)
	fontFaceRegex := regexp.MustCompile(`src: url\((?P<font>.*)\) `)
	fontFaceMatches := fontFaceRegex.FindStringSubmatch(face)
	if len(fontFaceMatches) != 2 {
		log.Printf("CPU %d: Failed to find url for font %s %s", idx, job.Name, job.Variant)
		log.Printf("CPU %d: %s", idx, face)
		return "", fmt.Errorf("failed to find font url")
	}

	fontIndex := fontFaceRegex.SubexpIndex("font")
	fontUrl := fontFaceMatches[fontIndex]

	return fontUrl, nil
}

func fetchCss(idx int, job FontDownloadJob, userAgent string) ([]byte, error) {
	log.Printf("CPU %d: Download font %s %s", idx, job.Name, job.Variant)
	query := ""
	if job.Variant == "regular" {
		query += "ital,wght@0,400"
	} else if job.Variant == "italic" {
		query += "ital,wght@1,400"
	} else if strings.HasSuffix(job.Variant, "italic") {
		query += "ital,wght@1," + strings.TrimSuffix(job.Variant, "italic")
	} else {
		query += "ital,wght@0," + job.Variant
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:%s", url.QueryEscape(job.Name), url.QueryEscape(query)), nil)
	if err != nil {
		log.Printf("CPU %d: Failed to create request for font %s %s", idx, job.Name, job.Variant)
		return []byte{}, err
	}

	req.Header.Add("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("CPU %d: Failed to get data for font %s %s", idx, job.Name, job.Variant)
		return []byte{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("CPU %d: Failed to get data for font %s %s", idx, job.Name, job.Variant)
		return []byte{}, err
	}

	fontCss, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("CPU %d: Failed to read body for font %s %s", idx, job.Name, job.Variant)
		return []byte{}, err
	}

	return fontCss, err
}

func Sync(configuration *config.Configuration) error {
	log.Println("Grab font list")
	fonts, err := downloadFontList(configuration.ApiKey)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	fontChannel := make(chan []FontDownloadJob)
	for i := 0; i < runtime.NumCPU(); i++ {
		go saveFontFile(configuration, fontChannel, wg, i)
	}

	for _, font := range fonts {
		variants := font.Variants
		name := font.Family
		category := font.Category
		var jobs []FontDownloadJob
		for _, variant := range variants {
			jobs = append(jobs, FontDownloadJob{
				Category: category,
				Name:     name,
				Variant:  variant,
			})
		}

		fontChannel <- jobs
	}

	close(fontChannel)
	wg.Wait()

	return nil
}
