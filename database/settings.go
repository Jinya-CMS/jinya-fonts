package database

import (
	"strings"
)

func CreateSettingsIfNotExists() (*JinyaFontsSettings, error) {
	settings, err := dbMap.SelectOneType[JinyaFontsSettings]("select * from settings limit 1")
	if err == nil {
		return settings, nil
	}

	settings = &JinyaFontsSettings{
		FilterByNameDb: "",
		SyncEnabled:    true,
		SyncInterval:   "0 0 1 * *",
	}

	err = dbMap.Insert(settings)

	return settings, err
}

func GetSettings() (*JinyaFontsSettings, error) {
	settings, err := dbMap.SelectOneType[JinyaFontsSettings]("select * from settings limit 1")

	return settings, err
}

func UpdateSettings(settings *JinyaFontsSettings) error {
	_, err := GetDbMap().Exec("update settings set filter_by_name = ?, sync_enabled = ?, sync_interval = ?", strings.Join(settings.FilterByName, ","), settings.SyncEnabled, settings.SyncInterval)

	return err
}
