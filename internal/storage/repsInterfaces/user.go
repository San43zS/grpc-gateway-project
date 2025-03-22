package repsInterfaces

import (
	"context"
	"grpc-gateway-project/internal/models"
)

//go:generate mockery --name User --output=./../testStorage/mocks
type User interface {
	Create(ctx context.Context, usr *models.User) (*models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	AddSubscribeEvent(ctx context.Context, ids *models.SubscribeEvent) error
	StealSubscribeEvent(ctx context.Context, ids *models.SubscribeEvent) error
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
}
