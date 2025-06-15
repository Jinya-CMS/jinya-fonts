package database

func ClearGoogleFonts() error {
	_, err := dbMap.Exec("delete from font where google_font = true")

	return err
}
