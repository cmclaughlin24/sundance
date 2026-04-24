package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBOpts struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Connect(opts ...func(*MongoDBOpts)) (*mongo.Client, error) {
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

func WithPort(port int) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Port = port
	}
}

func WithHost(host string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Host = host
	}
}

func WithUsername(username string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Username = username
	}
}

func WithPassword(password string) func(*MongoDBOpts) {
	return func(o *MongoDBOpts) {
		o.Password = password
	}
}

func Bootstrap(client *mongo.Client, logger *log.Logger) *ports.Repository {
	db := client.Database("tenants")

	return &ports.Repository{
		Database:    NewMongoDBDatabase(client, db),
		DataSources: newMongoDBDataSourcesRepository(db, logger),
		Tenants:     newMongoDBTenantsRepository(db, logger),
	}
}
