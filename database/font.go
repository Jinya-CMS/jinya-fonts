package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type Webfont struct {
	Name        string          `json:"name" bson:"name,omitempty"`
	Fonts       map[string]File `json:"fonts" bson:"fonts,omitempty"`
	Description string          `json:"description" bson:"description,omitempty"`
	Designers   []Designer      `json:"designers" bson:"designers,omitempty"`
	License     string          `json:"license" bson:"license,omitempty"`
	Category    string          `json:"category" bson:"category,omitempty"`
	GoogleFont  bool            `json:"googleFont" bson:"googleFont,omitempty"`
}

func GetAllFonts() ([]Webfont, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	cursor, err := fontsCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var fonts []Webfont
	err = cursor.All(ctx, &fonts)
	if err != nil {
		return nil, err
	}

	return fonts, nil
}

func GetGoogleFonts() ([]Webfont, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	cursor, err := fontsCollection.Find(ctx, bson.D{{"googleFont", true}})
	if err != nil {
		return nil, err
	}

	var fonts []Webfont
	err = cursor.All(ctx, &fonts)
	if err != nil {
		return nil, err
	}

	return fonts, nil
}

func GetCustomFonts() ([]Webfont, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	cursor, err := fontsCollection.Find(ctx, bson.D{{"googleFont", false}})
	if err != nil {
		return nil, err
	}

	var fonts []Webfont
	err = cursor.All(ctx, &fonts)
	if err != nil {
		return nil, err
	}

	return fonts, nil
}

func GetFont(name string) (*Webfont, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)

	font := new(Webfont)
	err = fontsCollection.FindOne(ctx, bson.D{{"name", name}}).Decode(font)
	if err != nil {
		return nil, err
	}

	return font, nil
}

func CreateFont(webfont *Webfont) error {
	client, err := openConnection()
	if err != nil {
		return err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	count, err := fontsCollection.CountDocuments(ctx, bson.D{{"name", webfont.Name}})
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("font already exists")
	}

	_, err = fontsCollection.InsertOne(ctx, webfont)

	return err
}

func UpdateFont(webfont *Webfont) error {
	client, err := openConnection()
	if err != nil {
		return err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	_, err = fontsCollection.UpdateOne(ctx, bson.D{{"name", webfont.Name}}, webfont)

	return err
}

func DeleteFont(name string) error {
	client, err := openConnection()
	if err != nil {
		return err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	_, err = fontsCollection.DeleteOne(ctx, bson.D{{"name", name}})

	return err
}

func ClearGoogleFonts() {
	client, err := openConnection()
	if err != nil {
		return
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	fontsCollection := getFontsCollection(client)
	_, _ = fontsCollection.DeleteMany(ctx, bson.D{{"googleFont", true}})
}
