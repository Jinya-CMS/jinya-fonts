package database

type Webfont struct {
	Name        string     `json:"name" db:"name,primarykey"`
	Fonts       []File     `json:"fonts" db:"-"`
	Description string     `json:"description" db:"description"`
	Designers   []Designer `json:"designers" db:"-"`
	License     string     `json:"license" db:"license"`
	Category    string     `json:"category" db:"category"`
	GoogleFont  bool       `json:"googleFont" db:"google_font"`
}

type Designer struct {
	Name string `json:"name" db:"name"`
	Bio  string `json:"bio" db:"bio"`
	Font string `json:"-" db:"font"`
}

type File struct {
	Path   string `json:"path" db:"path,primarykey"`
	Weight string `json:"weight" db:"weight"`
	Style  string `json:"style" db:"style"`
	Type   string `json:"type" db:"type"`
	Font   string `json:"-" db:"font"`
}

type JinyaFontsSettings struct {
	FilterByName   []string `json:"filterByName" db:"-"`
	FilterByNameDb string   `json:"-" db:"filter_by_name"`
	SyncEnabled    bool     `json:"syncEnabled" db:"sync_enabled"`
	SyncInterval   string   `json:"syncInterval" db:"sync_interval"`
}
