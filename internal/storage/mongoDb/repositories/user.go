package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/internal/storage/mongoDb/mongoServices"
	"grpc-gateway-project/internal/storage/repsInterfaces"
)

var userNotFoundErr = errors.New("user not found")

type Users struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewUserRep(collection *mongo.Collection, client *mongo.Client) repsInterfaces.User {
	return &Users{
		collection: collection,
		client:     client,
	}
}

func (s *Users) Create(ctx context.Context, usr *models.User) (*models.User, error) {
	id, err := s.collection.InsertOne(ctx, usr)
	if err != nil {
		return nil, err
	}
	usr.Id = id.InsertedID.(primitive.ObjectID).Hex()
	return usr, nil
}

func (s *Users) Get(ctx context.Context, id string) (*models.User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var usr models.User
	err = s.collection.FindOne(ctx, bson.D{{"_id", userId}}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (s *Users) Delete(ctx context.Context, id string) error {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	session, err := mongoServices.StartTransaction(ctx, s.client)
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	filter := bson.M{
		"$or": []bson.M{
			{"subscriptions": id},
			{"subscribers": id},
		},
	}
	update := bson.M{"$pull": bson.M{"subscriptions": id, "subscribers": id}}

	_, err = s.collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}

	res, err := s.collection.DeleteOne(ctx, bson.D{{"_id", userId}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return userNotFoundErr
	}
	return session.CommitTransaction(ctx)
}

func (s *Users) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	var usr models.User
	err := s.collection.FindOne(ctx, bson.D{{"token", token}}).Decode(&usr)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}
