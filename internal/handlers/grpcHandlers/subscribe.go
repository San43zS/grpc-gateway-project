package grpcHandlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/proto/api/generate/desc"

	"context"
)

func (h *Handlers) SubscribeUser(ctx context.Context, usr *desc.UserRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to get user data")
	}

	if err := h.scenarios.SubscribeUser(ctx, &models.SubscribeEvent{
		SubscriberId: user.Id,
		ListenerId:   usr.GetId(),
	}); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *Handlers) UnsubscribeUser(ctx context.Context, usr *desc.UserRequest) (*emptypb.Empty, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to get user data")
	}

	if err := h.scenarios.UnsubscribeUser(ctx, &models.SubscribeEvent{
		SubscriberId: user.Id,
		ListenerId:   usr.GetId(),
	}); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}
