package mongoServices

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartTransaction(ctx context.Context, client *mongo.Client) (mongo.Session, error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	err = session.StartTransaction()
	if err != nil {
		session.EndSession(ctx)
		return nil, err
	}
	return session, nil
}
