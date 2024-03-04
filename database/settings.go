package database

import (
	"go.mongodb.org/mongo-driver/bson"
)

type JinyaFontsSettings struct {
	FilterByName []string `json:"filterByName" bson:"filterByName"`
	SyncEnabled  bool     `json:"syncEnabled" bson:"syncEnabled"`
	SyncInterval string   `json:"syncInterval" bson:"syncInterval"`
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
	_, err = settingsCollection.DeleteMany(ctx, bson.D{})

	if err != nil {
		return err
	}

	_, err = settingsCollection.InsertOne(ctx, settings)

	return err
}
