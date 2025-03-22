package grpcHandlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"grpc-gateway-project/internal/models"
	"grpc-gateway-project/internal/scenarios"
	"grpc-gateway-project/proto/api/generate/desc"

	"context"
)

type Handlers struct {
	scenarios scenarios.User
	desc.UnimplementedUserServiceServer
}

func New(scenarios scenarios.User) *Handlers {
	return &Handlers{
		scenarios: scenarios,
	}
}

func (h *Handlers) CreateUser(ctx context.Context, usr *desc.UserData) (*desc.UserAccessInfo, error) {
	newUser, err := h.scenarios.CreateUser(
		ctx,
		&models.User{
			Email: usr.GetEmail(),
			Name:  usr.GetName(),
			Age:   usr.GetAge(),
		},
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &desc.UserAccessInfo{
		Token: newUser.Token,
		Id:    newUser.Id,
	}, nil
}

func (h *Handlers) GetUser(ctx context.Context, usr *desc.UserRequest) (*desc.User, error) {
	user, err := h.scenarios.GetUser(ctx, usr.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &desc.User{
		Id: user.Id,
		User: &desc.UserData{
			Email: user.Email,
			Name:  user.Name,
			Age:   user.Age,
		},
	}, nil
}

func (h *Handlers) DeleteUser(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to get user data")
	}

	if err := h.scenarios.DeleteUser(ctx, user.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}
