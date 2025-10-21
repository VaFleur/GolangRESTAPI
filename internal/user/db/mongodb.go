package db

import (
	"REST_API_app/internal/apperror"
	"REST_API_app/internal/user"
	"REST_API_app/pkg/logging"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to: %v", err)
	}

	d.logger.Debug("convert inserted ID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert ObjectID to HEX. Probably ObjectID: %s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert HEX to ObjectID. HEX: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find user by ID: %s due to error: %v", id, err)
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user ID: %s from DB due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	filter := bson.M{}
	cursor, err := d.collection.Find(ctx, filter)
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error: %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents form cursor due to error: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ObjectID: %s", user.ID)
	}

	filter := bson.M{"_id": oid}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user due to: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal userBytes due to: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{"$set": updateUserObj}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user due to: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ObjectID: %s", id)
	}

	filter := bson.M{"_id": oid}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user ID: %s due to: %v", id, err)
	}

	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
