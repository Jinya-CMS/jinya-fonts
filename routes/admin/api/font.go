package api

import (
	"encoding/json"
	"jinya-fonts/database"
	"jinya-fonts/fontsync"
	"jinya-fonts/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type addFontData struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type updateFontData struct {
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type apiFont struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Designers   []database.Designer `json:"designers"`
	License     string              `json:"license"`
	Category    string              `json:"category"`
	Fonts       []database.File     `json:"fonts"`
	GoogleFont  bool                `json:"googleFont"`
}

func getWebFont(name string) (*database.Webfont, []database.Designer, []database.File, error) {
	font, err := database.Get[database.Webfont](name)
	if err != nil {
		return nil, nil, nil, err
	}

	designers, err := database.GetDesigners(name)
	if err != nil {
		return nil, nil, nil, err
	}

	files, err := storage.GetFontFiles(name)
	if err != nil {
		return nil, nil, nil, err
	}

	return font, designers, files, nil
}

func convertWebfontToApiFont(webfont *database.Webfont, designers []database.Designer, files []database.File) apiFont {
	return apiFont{
		Fonts:       files,
		Name:        webfont.Name,
		Description: webfont.Description,
		Designers:   designers,
		License:     webfont.License,
		Category:    webfont.Category,
		GoogleFont:  webfont.GoogleFont,
	}
}

func groupFilesAndDesignersByFont(dbFiles []database.File, dbDesigners []database.Designer) (map[string][]database.Designer, map[string][]database.File, error) {
	files := make(map[string][]database.File)
	for _, file := range dbFiles {
		files[file.Font] = append(files[file.Font], file)
	}

	designers := make(map[string][]database.Designer)
	for _, designer := range dbDesigners {
		designers[designer.Font] = append(designers[designer.Font], designer)
	}

	return designers, files, nil
}

func getWebFontsDetails(google bool) (map[string][]database.Designer, map[string][]database.File, error) {
	designers, err := database.Select[database.Designer]("select d.* from designer d join font f on f.name = d.font where f.google_font = $1", google)
	if err != nil {
		return nil, nil, err
	}

	files, err := database.Select[database.File]("select d.* from file d join font f on f.name = d.font where f.google_font = $1", google)
	if err != nil {
		return nil, nil, err
	}

	return groupFilesAndDesignersByFont(files, designers)
}

func getWebFontDetails(name string) ([]database.Designer, []database.File, error) {
	designers, err := database.Select[database.Designer]("select * from designer where font = $1", name)
	if err != nil {
		return nil, nil, err
	}

	files, err := database.Select[database.File]("select * from file where font = $1", name)
	if err != nil {
		return nil, nil, err
	}

	return designers, files, nil
}

func getAllWebFontDetails() (map[string][]database.Designer, map[string][]database.File, error) {
	designers, err := database.Select[database.Designer]("select * from designer")
	if err != nil {
		return nil, nil, err
	}

	files, err := database.Select[database.File]("select * from file")
	if err != nil {
		return nil, nil, err
	}

	return groupFilesAndDesignersByFont(files, designers)
}

func combineFontFilesAndDesigners(fonts []database.Webfont, designers map[string][]database.Designer, files map[string][]database.File) []apiFont {
	apiFonts := make([]apiFont, 0)

	for _, font := range fonts {
		des, ok := designers[font.Name]
		if !ok {
			des = []database.Designer{}
		}

		fs, ok := files[font.Name]
		if !ok {
			fs = []database.File{}
		}
		apiFonts = append(apiFonts, convertWebfontToApiFont(&font, des, fs))
	}

	return apiFonts
}

func getApiFontList(google bool) ([]apiFont, error) {
	fonts, err := database.Select[database.Webfont]("select * from font where google_font = $1", google)
	if err != nil {
		return nil, err
	}

	designers, files, err := getWebFontsDetails(google)
	if err != nil {
		return nil, err
	}

	return combineFontFilesAndDesigners(fonts, designers, files), nil
}

func getFullApiFontList() ([]apiFont, error) {
	fonts, err := database.Select[database.Webfont]("select * from font")
	if err != nil {
		return nil, err
	}

	designers, files, err := getAllWebFontDetails()
	if err != nil {
		return nil, err
	}

	return combineFontFilesAndDesigners(fonts, designers, files), nil
}

func GetAllFonts(w http.ResponseWriter, r *http.Request) {
	fonts, err := getFullApiFontList()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(fonts)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getGoogleFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := getApiFontList(true)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(availableFonts)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getCustomFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := getApiFontList(false)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(availableFonts)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func GetFontByName(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.Get[database.Webfont](fontName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	designers, files, err := getWebFontDetails(font.Name)

	err = json.NewEncoder(w).Encode(convertWebfontToApiFont(font, designers, files))
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func createFont(w http.ResponseWriter, r *http.Request) {
	body := new(addFontData)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	font := database.Webfont{
		Name:        body.Name,
		Description: body.Description,
		License:     body.License,
		Category:    body.Category,
		GoogleFont:  false,
	}

	err = database.Insert(font)
	if err != nil {
		http.Error(w, "Failed to create font", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(font)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateFont(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.Get[database.Webfont](fontName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := new(updateFontData)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	font.Description = body.Description
	font.License = body.License
	font.Category = body.Category

	rows, err := database.Update(font)
	if err != nil {
		http.Error(w, "Failed to update font", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFont(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	res, err := database.GetDbMap().Exec("delete from font where name = $1", fontName)
	if err != nil {
		http.Error(w, "Failed to delete font", http.StatusInternalServerError)
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to delete font", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func syncFonts(w http.ResponseWriter, _ *http.Request) {
	err := fontsync.Sync()
	if err != nil {
		http.Error(w, "Failed to sync fonts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
