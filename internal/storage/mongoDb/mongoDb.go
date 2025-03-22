package mongoDb

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grpc-gateway-project/internal/storage/mongoDb/repositories"
	"grpc-gateway-project/internal/storage/repsInterfaces"

	"context"
)

type MongoDb struct {
	config *ConnectionMongo
	client *mongo.Client
	db     *mongo.Database
	reps   *MongoReps
}

type MongoReps struct {
	user repsInterfaces.User
}

func New(ctx context.Context) (*MongoDb, error) {
	var mongoDBURL string
	config := NewConfig()

	mongoDBURL = fmt.Sprintf("mongodb://%s:%s", config.MongoHost, config.MongoPort)
	fmt.Printf("Try to connect: %s\n", mongoDBURL)
	clientOptions := options.Client().ApplyURI(mongoDBURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	mongoDb := &MongoDb{
		client: client,
		config: config,
		db:     client.Database(config.MongoDBName),
	}
	mongoDb.reps = mongoDb.ConfigureReps()

	return mongoDb, nil
}

func (s *MongoDb) ConfigureReps() *MongoReps {
	return &MongoReps{
		user: repositories.NewUserRep(s.db.Collection("users"), s.client),
	}
}

func (s *MongoDb) User() repsInterfaces.User { return s.reps.user }
