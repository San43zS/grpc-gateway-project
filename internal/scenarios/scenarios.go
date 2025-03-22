package scenarios

import (
	"context"
	"errors"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/internal/services"
	"grpc-gateway-project/internal/storage"
)

var (
	userNotFoundErr = errors.New("user not found")
	emailEmptyErr   = errors.New("email is empty")
)

type User interface {
	CreateUser(ctx context.Context, usr *models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	SubscribeUser(ctx context.Context, ids *models.SubscribeEvent) error
	UnsubscribeUser(ctx context.Context, ids *models.SubscribeEvent) error
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
}

type Scenarios struct {
	storage storage.Storage
}

func New(storage storage.Storage) User {
	return &Scenarios{
		storage: storage,
	}
}

func (s *Scenarios) CreateUser(ctx context.Context, usr *models.User) (*models.User, error) {
	token, err := services.CreateToken(usr.Email)
	if err != nil {
		return nil, err
	}
	usr.Token = token
	usr.Subscribers = []string{}
	usr.Subscriptions = []string{}

	if usr.Email == "" {
		return nil, emailEmptyErr
	}

	newUser, err := s.storage.User().Create(ctx, usr)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *Scenarios) GetUser(ctx context.Context, id string) (*models.User, error) {
	usr, err := s.storage.User().Get(ctx, id)
	if err != nil {
		return nil, userNotFoundErr
	}

	return usr, nil
}

func (s *Scenarios) DeleteUser(ctx context.Context, id string) error {
	if err := s.storage.User().Delete(ctx, id); err != nil {
		return userNotFoundErr
	}

	return nil
}

func (s *Scenarios) SubscribeUser(ctx context.Context, ids *models.SubscribeEvent) error {
	return s.storage.User().AddSubscribeEvent(ctx, ids)
}

func (s *Scenarios) UnsubscribeUser(ctx context.Context, ids *models.SubscribeEvent) error {
	return s.storage.User().StealSubscribeEvent(ctx, ids)
}

func (s *Scenarios) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	usr, err := s.storage.User().GetUserByToken(ctx, token)
	if err != nil {
		return nil, userNotFoundErr
	}

	return usr, nil
}
