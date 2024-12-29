package mongo

import (
	"context"

	"github.com/dro14/nuqta-service/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (m *Mongo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := m.users.InsertOne(ctx, user)
	return err
}

func (m *Mongo) ReadUser(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	filter := bson.M{"_id": id}
	err := m.users.FindOne(ctx, filter).Decode(user)
	return user, err
}

func (m *Mongo) UpdateUser(ctx context.Context, user *models.User) error {
	update := bson.M{"$set": user}
	filter := bson.M{"_id": user.ID}
	result, err := m.users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (m *Mongo) DeleteUser(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := m.users.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
