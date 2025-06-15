package database

import (
	"fmt"
)

func GetDesigners(name string) ([]Designer, error) {
	return Select[Designer]("select * from designer where font = $1", name)
}

func createDesigner(font *Webfont, designer Designer) (*Designer, error) {
	_, err := dbMap.Exec("insert into designer (name, bio, font) values ($1, $2, $3) on conflict do update set bio = $2", designer.Name, designer.Bio, font.Name)

	return &designer, err
}

func deleteDesigner(font *Webfont, designerName string) error {
	_, err := dbMap.Exec("delete from designer where font = $1 and name = $2", font.Name, designerName)

	return err
}

func CreateDesigner(name string, designer Designer) (*Designer, error) {
	font, err := Get[Webfont](name)
	if err != nil {
		return nil, err
	}

	if font.GoogleFont {
		return nil, fmt.Errorf("cannot add designers to a google font")
	}

	return createDesigner(font, designer)
}

func DeleteDesigner(name string, designerName string) error {
	font, err := Get[Webfont](name)
	if err != nil {
		return err
	}

	if font.GoogleFont {
		return fmt.Errorf("cannot remove designers from a google font")
	}

	return deleteDesigner(font, designerName)
}

func CreateGoogleDesigner(name string, designer Designer) (*Designer, error) {
	font, err := Get[Webfont](name)
	if err != nil {
		return nil, err
	}

	return createDesigner(font, designer)
}

func DeleteGoogleDesigner(name string, designerName string) error {
	font, err := Get[Webfont](name)
	if err != nil {
		return err
	}

	return deleteDesigner(font, designerName)
}
