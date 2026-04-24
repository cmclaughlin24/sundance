package mongodb

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Bootstrap(client *mongo.Client, logger *log.Logger) *ports.Repository {
	return &ports.Repository{}
}
