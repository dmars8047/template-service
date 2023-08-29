package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TemplateStore interface {
	GetAllTemplates(name string) ([]Template, error)
	GetTemplateById(templateId string) (*Template, error)
	GetTemplatebyName(templateName string) (*Template, error)
	CreateTemplate(template *Template) error
	DeleteTemplate(templateId string) error
}

type MongoTemplateStore struct {
	collection *mongo.Collection
}

func NewMongoTemplateStore(collection *mongo.Collection) *MongoTemplateStore {
	return &MongoTemplateStore{collection: collection}
}

func (store *MongoTemplateStore) GetAllTemplates(name string) ([]Template, error) {
	filter := bson.M{}

	if name != "" {
		filter["name"] = name
	}

	var results []Template = make([]Template, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := store.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (store *MongoTemplateStore) GetTemplateById(templateId string) (*Template, error) {
	filter := bson.M{"_id": templateId}

	var result Template

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

func (store *MongoTemplateStore) GetTemplatebyName(templateName string) (*Template, error) {
	filter := bson.M{"name": templateName}

	var result Template

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

func (store *MongoTemplateStore) CreateTemplate(template *Template) error {
	// Make sure a template with the same name doesn't already exist
	_, err := store.GetTemplatebyName(template.Name)

	if err == nil {
		return ErrNameConflict
	}

	if err == ErrNotFound {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := store.collection.InsertOne(ctx, bson.M{
			"_id":     template.Id,
			"name":    template.Name,
			"content": template.Content,
			"format":  template.Format,
			"tokens":  template.Tokens,
		})

		if err != nil {
			return err
		}

		return nil
	} else {
		return err
	}
}

func (store *MongoTemplateStore) DeleteTemplate(templateId string) error {
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
