package api

import (
	"context"
	"time"

	tmpl "github.com/dmars8047/marshall-labs-common/templates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailTemplateStore interface {
	GetEmailTemplate(templateId string) (*tmpl.EmailTemplate, error)
	CreateEmailTemplate(template *tmpl.EmailTemplate) error
	DeleteEmailTemplate(templateId string) error
}

type MongoEmailTemplateStore struct {
	collection *mongo.Collection
}

func NewMongoEmailTemplateStore(collection *mongo.Collection) *MongoEmailTemplateStore {
	return &MongoEmailTemplateStore{collection: collection}
}

func (store *MongoEmailTemplateStore) GetEmailTemplate(templateId string) (*tmpl.EmailTemplate, error) {
	filter := bson.M{"_id": templateId}

	var result tmpl.EmailTemplate

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := store.collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &result, nil
}

func (store *MongoEmailTemplateStore) GetEmailTemplateByName(templateName string) (*tmpl.EmailTemplate, error) {
	filter := bson.M{"name": templateName}

	var result tmpl.EmailTemplate

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := store.collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &result, nil
}

func (store *MongoEmailTemplateStore) CreateEmailTemplate(template *tmpl.EmailTemplate) error {
	// Make sure a template with the same name doesn't already exist
	_, err := store.GetEmailTemplateByName(template.Name)

	if err == nil {
		return ErrNameConflict
	}

	if err == ErrNotFound {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := store.collection.InsertOne(ctx, bson.M{
			"_id":                template.Id,
			"name":               template.Name,
			"html_content":       template.HtmlContent,
			"plain_text_content": template.PlainTextContent,
			"subject":            template.Subject,
			"tokens":             template.Tokens,
		})

		if err != nil {
			return err
		}

		return nil
	} else {
		return err
	}
}

func (store *MongoEmailTemplateStore) DeleteEmailTemplate(templateId string) error {
	filter := bson.M{"_id": templateId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := store.collection.DeleteOne(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}

		return err
	}

	return nil
}
