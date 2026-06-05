package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBOpts struct {
	URI          string `json:"uri" yaml:"uri"`
	DatabaseName string `json:"database_name" yaml:"databaseName"`
}

func ConnectMongoDB(opts ...func(*MongoDBOpts)) (*mongo.Client, error) {
	var o MongoDBOpts

	for _, opt := range opts {
		opt(&o)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(o.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(clientOpts)

	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

func MongoDBWithURI(uri string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.URI = uri
	}
}
