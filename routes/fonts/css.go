package fonts

import (
	"bytes"
	"encoding/csv"
	"jinya-fonts/database"
	"jinya-fonts/storage"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"text/template"
)

var tmpl = `
@font-face {
    font-family: '{{.Name}}';
    font-style: {{.Style}};
    font-weight: {{.Weight}};
    src: url('{{.Path}}') format('{{.Type}}');
    font-display: {{.FontDisplay}};
}
`

type templateData struct {
	*database.File
	*database.Webfont
	FontDisplay string
}

type family struct {
	Name   string
	Italic bool
	Weight string
}

func convertFamilyToTemplateData(fam family, display string) []templateData {
	webfont, err := database.Get[database.Webfont](fam.Name)
	if err != nil {
		return []templateData{}
	}

	files, err := storage.GetFontFiles(webfont.Name)
	if err != nil {
		return []templateData{}
	}

	var data []templateData
	for _, metadata := range files {
		if fam.Italic && metadata.Style != "italic" {
			continue
		}
		if metadata.Weight == fam.Weight {
			data = append(data, templateData{
				File:        &metadata,
				Webfont:     webfont,
				FontDisplay: display,
			})
		}
	}

	return data
}

func getFamiliesFromModifiers(allModifiers [][]string, name string) []family {
	italicIdx := -1
	weightIdx := -1
	headerRow := allModifiers[0]
	for i := range headerRow {
		if headerRow[i] == "ital" {
			italicIdx = i
		} else if headerRow[i] == "wght" {
			weightIdx = i
		}
	}

	var families []family

	for _, entry := range allModifiers[1:] {
		italic := false
		weight := "regular"
		if italicIdx != -1 {
			italic = entry[italicIdx] == "1"
		}

		if weightIdx != -1 {
			weight = entry[weightIdx]
		}

		families = append(families, family{name, italic, weight})
	}

	return families
}

func convertTemplateDataToCss(item templateData) (string, error) {
	parsedTmpl, err := template.New("css").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = parsedTmpl.Execute(&buf, item)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getCss2(w http.ResponseWriter, r *http.Request) {
	unescapedQuery, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	querySplitIntoFamilies := strings.Split(unescapedQuery, "&")

	var families []family
	display := "swap"

	for _, query := range querySplitIntoFamilies {
		if strings.HasPrefix(query, "family=") {
			trimmed := strings.TrimPrefix(query, "family=")
			splitFamily := strings.Split(trimmed, ":")
			if len(splitFamily) == 1 {
				families = append(families, family{splitFamily[0], false, "400"})
			} else if len(splitFamily) == 2 && slices.Contains(splitFamily, "all") {
				font, err := database.Get[database.Webfont](splitFamily[0])
				if err != nil {
					http.NotFound(w, r)
					return
				}

				files, err := storage.GetFontFiles(font.Name)
				if err != nil {
					http.NotFound(w, r)
					return
				}

				for _, file := range files {
					families = append(families, family{
						Name:   font.Name,
						Italic: file.Style == "italic",
						Weight: file.Weight,
					})
				}
			} else if len(splitFamily) >= 2 {
				weightAndItalic := strings.ReplaceAll(strings.ReplaceAll(splitFamily[1], ";", "\n"), "@", "\n")
				reader := csv.NewReader(bytes.NewBufferString(weightAndItalic))
				reader.Comma = ','
				allModifiers, err := reader.ReadAll()
				if err != nil {
					http.NotFound(w, r)
					return
				}

				families = append(families, getFamiliesFromModifiers(allModifiers, splitFamily[0])...)
			}
		} else if strings.HasPrefix(query, "display=") {
			display = strings.TrimPrefix(query, "display=")
		}
	}

	var data []templateData
	for _, fam := range families {
		data = append(data, convertFamilyToTemplateData(fam, display)...)
	}

	css := ""
	for _, item := range data {
		data, err := convertTemplateDataToCss(item)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		css += data
	}

	w.Header().Add("Content-Type", "text/css")
	_, _ = w.Write([]byte(css))
}
