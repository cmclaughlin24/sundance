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
	forms    *database.MongoDBRepository[formDocument]
	versions *database.MongoDBRepository[any]
}

func NewMongoDBFormsRepository(db *mongo.Database, logger *log.Logger) ports.FormsRepository {
	formsRepository := database.NewMongoDBRepository[formDocument](
		db.Collection("forms"),
		logger,
	)
	versionsRepository := database.NewMongoDBRepository[any](
		db.Collection("versions"),
		logger,
	)

	return &mongoDBFormsRepository{
		forms:    formsRepository,
		versions: versionsRepository,
	}
}

func (r *mongoDBFormsRepository) Find(ctx context.Context) ([]*domain.Form, error) {
	documents, err := r.forms.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	forms := make([]*domain.Form, 0, len(documents))
	for _, document := range documents {
		forms = append(forms, fromFormDocument(&document))
	}

	return forms, nil
}

func (r *mongoDBFormsRepository) FindById(ctx context.Context, formID domain.FormID) (*domain.Form, error) {
	document, err := r.forms.FindOne(ctx, bson.M{"_id": formID})

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
		return r.forms.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromFormDocument(&result), nil
}

func (r *mongoDBFormsRepository) FindVersions(context.Context, domain.FormID) ([]*domain.Version, error) {
	return nil, nil
}

func (r *mongoDBFormsRepository) FindVersion(context.Context, domain.FormID, domain.VersionID) (*domain.Version, error) {
	return nil, nil
}

func (r *mongoDBFormsRepository) FindNextVersionNumber(context.Context, domain.FormID) (int, error) {
	return 0, nil
}

func (r *mongoDBFormsRepository) UpsertVersion(context.Context, *domain.Version) (*domain.Version, error) {
	return nil, nil
}
