package api

import (
	"context"
	"time"

	"github.com/dmars8047/templates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailTemplateStore interface {
	getEmailTemplate(templateId string) (*templates.EmailTemplate, error)
	createEmailTemplate(template *templates.EmailTemplate) error
	deleteEmailTemplate(templateId string) error
}

type MongoEmailTemplateStore struct {
	collection *mongo.Collection
}

func NewMongoEmailTemplateStore(collection *mongo.Collection) *MongoEmailTemplateStore {
	return &MongoEmailTemplateStore{collection: collection}
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
