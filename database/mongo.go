package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"jinya-fonts/config"
)

func openConnection() (*mongo.Client, error) {
	ctx, cancelFunc := getContext()
	defer cancelFunc()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.LoadedConfiguration.MongoUrl))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func closeConnection(client *mongo.Client) {
	ctx, cancelFunc := getContext()
	defer cancelFunc()

	_ = client.Disconnect(ctx)
}

func getDatabase(client *mongo.Client) *mongo.Database {
	return client.Database(config.LoadedConfiguration.MongoDatabase)
}

func getFontsCollection(client *mongo.Client) *mongo.Collection {
	return getDatabase(client).Collection("fonts")
}

func getSettingsCollection(client *mongo.Client) *mongo.Collection {
	return getDatabase(client).Collection("settings")
}

func CheckMongo() bool {
	ctx, cancelFunc := getContext()
	defer cancelFunc()

	client, err := openConnection()
	if err != nil {
		return false
	}

	defer closeConnection(client)

	return client.Ping(ctx, readpref.Primary()) == nil
}
