package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JinyaFontsSettings struct {
	FilterByName []string `json:"filterByName" bson:"filterByName"`
}

func GetSettings() (*JinyaFontsSettings, error) {
	client, err := openConnection()
	if err != nil {
		return nil, err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	settingsCollection := getSettingsCollection(client)

	settings := new(JinyaFontsSettings)
	err = settingsCollection.FindOne(ctx, bson.D{}).Decode(settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func UpdateSettings(settings *JinyaFontsSettings) error {
	client, err := openConnection()
	if err != nil {
		return err
	}

	defer closeConnection(client)

	ctx, cancelFunc := getContext()
	defer cancelFunc()

	settingsCollection := getSettingsCollection(client)

	_, err = settingsCollection.UpdateOne(ctx, bson.D{}, settings, options.Update().SetUpsert(true))

	return err
}
