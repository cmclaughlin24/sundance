package mongodb

import (
	"context"
	"log"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoDBFormsRepository struct {
	base *database.MongoDBRepository[formDocument]
}

func newMongoDBFormsRepository(db *mongo.Database, logger *log.Logger) (ports.FormsRepository, error) {
	base := database.NewMongoDBRepository[formDocument](
		db.Collection("forms"),
		logger,
	)

	return &mongoDBFormsRepository{base}, nil
}

func (r *mongoDBFormsRepository) Find(ctx context.Context, f *ports.FormFilters) ([]*domain.Form, error) {
	filter := bson.M{}

	if f != nil && f.TenantID != "" {
		filter["tenant_id"] = f.TenantID
	}

	documents, err := r.base.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	forms := make([]*domain.Form, 0, len(documents))
	for _, document := range documents {
		forms = append(forms, fromFormDocument(&document))
	}

	return forms, nil
}

func (r *mongoDBFormsRepository) FindByID(ctx context.Context, formID domain.FormID) (*domain.Form, error) {
	document, err := r.base.FindOne(ctx, bson.M{"_id": formID})

	if err != nil {
		return nil, err
	}

	return fromFormDocument(&document), nil
}

func (r *mongoDBFormsRepository) Upsert(ctx context.Context, f *domain.Form) (*domain.Form, error) {
	doc := toFormDocument(f)
	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result formDocument
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.base.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromFormDocument(&result), nil
}

func (r *mongoDBFormsRepository) Delete(ctx context.Context, id domain.FormID) error {
	return r.base.Remove(ctx, bson.M{"_id": id})
}
