package mongodb

import (
	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBDatabase struct {
	db *mongo.Database
}

func NewMongoDBDatabase(db *mongo.Database) database.Database {
	return &MongoDBDatabase{
		db: db,
	}
}
