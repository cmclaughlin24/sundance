package mongodb

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *log.Logger) *ports.Repository {
	db := client.Database("submissions")

	return &ports.Repository{
		Database:    database.NewMongoDBDatabase(client, db),
		Submissions: newMongoDBSubmissionsRepository(db, logger),
	}
}
