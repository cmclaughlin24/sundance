package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBOpts struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ConnectMongoDB(opts ...func(*MongoDBOpts)) (*mongo.Client, error) {
	var o MongoDBOpts

	for _, opt := range opts {
		opt(&o)
	}

	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
	)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(clientOpts)

	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

func MongoDBWithPort(port int) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Port = port
	}
}

func MongoDBWithHost(host string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Host = host
	}
}

func MongoDBWithUsername(username string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Username = username
	}
}

func MongoDBWithPassword(password string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Password = password
	}
}
