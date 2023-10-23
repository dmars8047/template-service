package api

import (
	"context"
	"time"

	"github.com/dmars8047/template-service/pkg/templates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EmailTemplateHandler is a struct that handles get requests for email templates
type EmailTemplateGetter interface {
	getAllEmailTemplates() ([]templates.EmailTemplate, error)
	getEmailTemplate(templateId string) (*templates.EmailTemplate, error)
}

// EmailTemplateHandler is a struct that handles create requests for email templates
type EmailTemplateCreator interface {
	createEmailTemplate(template *templates.EmailTemplate) error
}

// EmailTemplateHandler is a struct that handles delete requests for email templates
type EmailTemplateDeleter interface {
	deleteEmailTemplate(templateId string) error
}

// MongoEmailTemplateStore is a struct that implements the EmailTemplateGetter, EmailTemplateCreator, and EmailTemplateDeleter interfaces
type MongoEmailTemplateStore struct {
	collection *mongo.Collection
}

func NewMongoEmailTemplateStore(collection *mongo.Collection) *MongoEmailTemplateStore {
	return &MongoEmailTemplateStore{collection: collection}
}

func (store *MongoEmailTemplateStore) getAllEmailTemplates() ([]templates.EmailTemplate, error) {
	var results []templates.EmailTemplate

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := store.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (store *MongoEmailTemplateStore) getEmailTemplate(templateId string) (*templates.EmailTemplate, error) {
	filter := bson.M{"_id": templateId}

	var result templates.EmailTemplate

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

func (store *MongoEmailTemplateStore) getEmailTemplateByName(templateName string) (*templates.EmailTemplate, error) {
	filter := bson.M{"name": templateName}

	var result templates.EmailTemplate

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

func (store *MongoEmailTemplateStore) createEmailTemplate(template *templates.EmailTemplate) error {
	// Make sure a template with the same name doesn't already exist
	_, err := store.getEmailTemplateByName(template.Name)

	if err == nil {
		return ErrNameConflict
	}

	if err == ErrNotFound {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := store.collection.InsertOne(ctx, bson.M{
			"_id":              template.Id,
			"name":             template.Name,
			"htmlContent":      template.HtmlContent,
			"plainTextContent": template.PlainTextContent,
			"subject":          template.Subject,
			"tokens":           template.Tokens,
		})

		if err != nil {
			return err
		}

		return nil
	} else {
		return err
	}
}

func (store *MongoEmailTemplateStore) deleteEmailTemplate(templateId string) error {
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
