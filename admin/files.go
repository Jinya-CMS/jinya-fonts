package admin

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"os"
	"sort"
	"strings"
)

func FilesIndex(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"You cannot edit the files of a synced font", fontName, []meta.FontFileMeta{}})
		return
	}

	sort.Slice(font.Fonts, func(i, j int) bool {
		first := font.Fonts[i]
		second := font.Fonts[j]

		return strings.ToLower(first.Path) < strings.ToLower(second.Path)
	})
	err = RenderAdmin(w, "files/index", struct {
		Message  string
		FontName string
		Files    []meta.FontFileMeta
	}{"", font.Name, font.Fonts})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "files/delete", struct {
			Message  string
			FontName string
			Path     string
		}{"Font does not exist", fontName, ""})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Path     string
		}{"Font does not exist", fontName, ""})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Path     string
		}{"You cannot edit the files of a synced font", fontName, ""})
		return
	}

	path := r.URL.Query().Get("path")
	var file *meta.FontFileMeta
	for _, item := range font.Fonts {
		if item.Path == path {
			file = &item
			break
		}
	}

	if file == nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Path     string
		}{"File doesn't exist in font", fontName, ""})
		return
	}

	if r.Method == http.MethodGet {
		RenderAdmin(w, "files/delete", struct {
			Message  string
			FontName string
			Path     string
		}{"", fontName, path})
	} else if r.Method == http.MethodPost {
		var files []meta.FontFileMeta
		for _, item := range font.Fonts {
			if item.Path != path {
				files = append(files, item)
			}
		}

		font.Fonts = files

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "files/delete", struct {
				Message  string
				FontName string
				Path     string
			}{"Failed to remove file from font", fontName, path})
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "files/delete", struct {
				Message  string
				FontName string
				Path     string
			}{"Failed to remove file from font", fontName, path})
			return
		}

		http.Redirect(w, r, "/admin/files?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddFile(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"You cannot edit the files of a synced font", fontName, []meta.FontFileMeta{}})
		return
	}

	if r.Method == http.MethodGet {
		RenderAdmin(w, "files/add", struct {
			Message  string
			FontName string
			Subset   string
			Weight   string
			Style    string
		}{"", fontName, "", "", "normal"})
	} else if r.Method == http.MethodPost {
		err = r.ParseMultipartForm(10 * 1024 * 1024 * 1024)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, "", "", ""})
			return
		}

		fileStyle := r.FormValue("style")
		fileSubset := r.FormValue("subset")
		fileWeight := r.FormValue("weight")
		filename := fontName + "." + fileSubset + "." + fileWeight + fileStyle + ".woff2"

		fileExists := false
		for _, item := range font.Fonts {
			if strings.ToLower(item.Path) == strings.ToLower(filename) {
				fileExists = true
				break
			}
		}

		if fileExists {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"A file with the chosen properties is already present", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		fontFile, _, err := r.FormFile("file")
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		fileBuffer := bytes.NewBufferString("")
		_, err = io.Copy(fileBuffer, fontFile)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+fontName, 0775)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+"/"+filename, fileBuffer.Bytes(), 0775)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		font.Fonts = append(font.Fonts, meta.FontFileMeta{
			Path:     filename,
			Subset:   fileSubset,
			Weight:   fileWeight,
			Style:    fileStyle,
			FontName: fontName,
		})

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "files/add", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		http.Redirect(w, r, "/admin/files?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func EditFile(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	path := r.URL.Query().Get("path")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"You cannot edit the files of a synced font", fontName, []meta.FontFileMeta{}})
		return
	}
	var fontFile meta.FontFileMeta
	for _, item := range font.Fonts {
		if item.Path == path {
			fontFile = item
			break
		}
	}

	if r.Method == http.MethodGet {
		RenderAdmin(w, "files/edit", struct {
			Message  string
			FontName string
			Subset   string
			Weight   string
			Style    string
		}{"", fontName, fontFile.Subset, fontFile.Weight, fontFile.Style})
	} else if r.Method == http.MethodPost {
		err = r.ParseMultipartForm(10 * 1024 * 1024 * 1024)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to add file to font", fontName, fontFile.Subset, fontFile.Weight, fontFile.Style})
			return
		}

		fileStyle := r.FormValue("style")
		fileSubset := r.FormValue("subset")
		fileWeight := r.FormValue("weight")
		filename := fontName + "." + fileSubset + "." + fileWeight + fileStyle + ".woff2"

		fileExists := false
		for _, item := range font.Fonts {
			if strings.ToLower(item.Path) == strings.ToLower(filename) && strings.ToLower(item.Path) != strings.ToLower(path) {
				fileExists = true
				break
			}
		}

		if fileExists {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"A file with the chosen properties is already present", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		_ = os.Remove(config.LoadedConfiguration.FontFileFolder + "/" + fontName + "/" + path)

		fontFile, _, err := r.FormFile("file")
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		fileBuffer := bytes.NewBufferString("")
		_, err = io.Copy(fileBuffer, fontFile)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+fontName, 0775)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+"/"+filename, fileBuffer.Bytes(), 0775)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		for idx, item := range font.Fonts {
			if strings.ToLower(item.Path) == strings.ToLower(path) {
				font.Fonts[idx].Path = filename
				font.Fonts[idx].Style = fileStyle
				font.Fonts[idx].Weight = fileWeight
				font.Fonts[idx].Subset = fileSubset
			}
		}

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "files/edit", struct {
				Message  string
				FontName string
				Subset   string
				Weight   string
				Style    string
			}{"Failed to update file in font", fontName, fileSubset, fileWeight, fileStyle})
			return
		}

		http.Redirect(w, r, "/admin/files?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
