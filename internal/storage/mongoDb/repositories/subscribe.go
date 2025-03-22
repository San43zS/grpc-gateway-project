package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/internal/storage/mongoDb/mongoServices"
)

type editSubscribeFunc func(ctx context.Context, ids *models.SubscribeEvent) error
type editSubscribeCreator func(ids *subscribeEvent) *subscribeEditor

type subscribeEvent struct {
	SubscriberId primitive.ObjectID
	ListenerId   primitive.ObjectID
}

type subscribeEditor struct {
	filter bson.M
	update bson.M
}

func (s *Users) AddSubscribeEvent(ctx context.Context, ids *models.SubscribeEvent) error {
	return s.editSubscribeEvent(s.createAddSubscribe, s.createAddSubscriber)(ctx, ids)
}

func (s *Users) StealSubscribeEvent(ctx context.Context, ids *models.SubscribeEvent) error {
	return s.editSubscribeEvent(s.createStealSubscribe, s.createStealSubscriber)(ctx, ids)
}

func (s *Users) editSubscribeEvent(subscriberEditFunc editSubscribeCreator, listenerEditFunc editSubscribeCreator) editSubscribeFunc {
	return func(ctx context.Context, ids *models.SubscribeEvent) error {
		subscriberId, err := primitive.ObjectIDFromHex(ids.SubscriberId)
		if err != nil {
			return err
		}
		listenerId, err := primitive.ObjectIDFromHex(ids.ListenerId)
		if err != nil {
			return err
		}

		dbIds := &subscribeEvent{
			SubscriberId: subscriberId,
			ListenerId:   listenerId,
		}

		session, err := mongoServices.StartTransaction(ctx, s.client)
		if err != nil {
			return err
		}
		defer session.EndSession(ctx)

		if err = s.editSubscribes(ctx, subscriberEditFunc(dbIds)); err != nil {
			return err
		}
		if err = s.editSubscribes(ctx, listenerEditFunc(dbIds)); err != nil {
			return err
		}
		return session.CommitTransaction(ctx)
	}
}

func (s *Users) createAddSubscribe(ids *subscribeEvent) *subscribeEditor {
	return &subscribeEditor{
		filter: bson.M{"_id": ids.SubscriberId},
		update: bson.M{"$addToSet": bson.M{"subscriptions": ids.ListenerId.Hex()}},
	}
}

func (s *Users) createAddSubscriber(ids *subscribeEvent) *subscribeEditor {
	return &subscribeEditor{
		filter: bson.M{"_id": ids.ListenerId},
		update: bson.M{"$addToSet": bson.M{"subscribers": ids.SubscriberId.Hex()}},
	}
}

func (s *Users) createStealSubscribe(ids *subscribeEvent) *subscribeEditor {
	return &subscribeEditor{
		filter: bson.M{"_id": ids.SubscriberId},
		update: bson.M{"$pull": bson.M{"subscriptions": ids.ListenerId.Hex()}},
	}
}

func (s *Users) createStealSubscriber(ids *subscribeEvent) *subscribeEditor {
	return &subscribeEditor{
		filter: bson.M{"_id": ids.ListenerId},
		update: bson.M{"$pull": bson.M{"subscribers": ids.SubscriberId.Hex()}},
	}
}

func (s *Users) editSubscribes(ctx context.Context, editor *subscribeEditor) error {
	res, err := s.forceEditSubscribes(ctx, editor)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return userNotFoundErr
	}
	return nil
}

func (s *Users) forceEditSubscribes(ctx context.Context, editor *subscribeEditor) (*mongo.UpdateResult, error) {
	res, err := s.collection.UpdateOne(ctx, editor.filter, editor.update)
	if err != nil {
		return nil, err
	}
	return res, nil
}
