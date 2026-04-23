package mongodb

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongodbBaseRepository struct {
	collection *mongo.Collection
	logger     *log.Logger
}
