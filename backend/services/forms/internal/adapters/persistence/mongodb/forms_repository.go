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
	versions *database.MongoDBRepository[versionDocument]
}

func newMongoDBFormsRepository(db *mongo.Database, logger *log.Logger) ports.FormsRepository {
	formsRepository := database.NewMongoDBRepository[formDocument](
		db.Collection("forms"),
		logger,
	)
	versionsRepository := database.NewMongoDBRepository[versionDocument](
		db.Collection("versions"),
		logger,
	)

	return &mongoDBFormsRepository{
		forms:    formsRepository,
		versions: versionsRepository,
	}
}

func (r *mongoDBFormsRepository) Find(ctx context.Context, f *ports.FormFilters) ([]*domain.Form, error) {
	var filter bson.M

	if f != nil && f.TenantID != "" {
		filter["tenant_id"] = f.TenantID
	}

	documents, err := r.forms.Find(ctx, filter)

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

func (r *mongoDBFormsRepository) FindVersions(ctx context.Context, formID domain.FormID) ([]*domain.Version, error) {
	documents, err := r.versions.Find(ctx, bson.M{"form_id": formID})

	if err != nil {
		return nil, err
	}

	versions := make([]*domain.Version, 0, len(documents))
	for _, document := range documents {
		v, err := fromVersionDocument(&document)

		if err != nil {
			return nil, err
		}

		versions = append(versions, v)
	}

	return versions, nil
}

func (r *mongoDBFormsRepository) FindVersion(ctx context.Context, formID domain.FormID, versionID domain.VersionID) (*domain.Version, error) {
	document, err := r.versions.FindOne(ctx, bson.M{"form_id": formID, "_id": versionID})

	if err != nil {
		return nil, err
	}

	return fromVersionDocument(&document)
}

func (r *mongoDBFormsRepository) FindNextVersionNumber(ctx context.Context, formID domain.FormID) (int, error) {
	filter := bson.M{"form_id": formID}
	opts := options.Find().SetSort(bson.M{"version": -1}).SetLimit(1).SetProjection(bson.M{"version": 1})

	var lv int
	err := mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		cursor, err := r.versions.Collection().Find(sctx, filter, opts)

		if err != nil {
			return err
		}

		defer cursor.Close(sctx)

		var documents []versionDocument
		if err := cursor.All(sctx, &documents); err != nil {
			return err
		}

		if len(documents) == 0 {
			lv = 1
			return nil
		}

		lv = documents[0].Version + 1
		return nil
	})

	if err != nil {
		return 0, err
	}

	return lv, nil
}

func (r *mongoDBFormsRepository) UpsertVersion(ctx context.Context, v *domain.Version) (*domain.Version, error) {
	doc, err := toVersionDocument(v)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result versionDocument
	err = mongo.WithSession(ctx, mongo.SessionFromContext(ctx), func(sctx context.Context) error {
		return r.versions.Collection().FindOneAndUpdate(sctx, filter, update, opts).Decode(&result)
	})

	if err != nil {
		return nil, err
	}

	return fromVersionDocument(&result)
}
