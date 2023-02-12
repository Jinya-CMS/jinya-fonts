package http

import (
	"bytes"
	"encoding/csv"
	"jinya-fonts/meta"
	"net/http"
	"strings"
	"text/template"
)

var tmpl = `
/* {{.Subset}} */
@font-face {
    font-family: '{{.Name}}';
    font-style: {{.Style}};
    font-weight: {{.Weight}};
    src: url('{{.Url}}') format('woff2');
	{{if ne .UnicodeRange ""}}
    unicode-range: {{.UnicodeRange}};
	{{end}}
    font-display: {{.FontDisplay}};
}
`

type templateData struct {
	Subset       string
	Name         string
	Style        string
	Url          string
	UnicodeRange string
	Weight       string
	FontDisplay  string
}

type family struct {
	Name   string
	Italic bool
	Weight string
}

func convertFamilyToTemplateData(fam family, display string) []templateData {
	fontFile, err := meta.LoadFontFileCache(fam.Name)
	if err != nil {
		return []templateData{}
	}

	var data []templateData
	for _, fontFamily := range fontFile.Fonts {
		if fam.Italic && fontFamily.Style != "italic" {
			continue
		}
		if fontFamily.Weight == fam.Weight {
			data = append(data, templateData{
				Subset:       fontFamily.Subset,
				Name:         fam.Name,
				Style:        fontFamily.Style,
				Url:          "/fonts/" + fam.Name + "/" + fontFamily.Path,
				UnicodeRange: fontFamily.UnicodeRange,
				Weight:       fam.Weight,
				FontDisplay:  display,
			})
		}
	}

	return data
}

func GetCss2(w http.ResponseWriter, r *http.Request) {
	unescapedQuery := strings.TrimPrefix(r.RequestURI, "/css2?")
	querySplitIntoFamilies := strings.Split(unescapedQuery, "&")

	var families []family
	display := "swap"

	for _, query := range querySplitIntoFamilies {
		if strings.HasPrefix(query, "family=") {
			trimmed := strings.TrimPrefix(query, "family=")
			splitFamily := strings.Split(trimmed, ":")
			if len(splitFamily) == 1 {
				families = append(families, family{splitFamily[0], false, "400"})
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
	w.Write([]byte(css))
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
