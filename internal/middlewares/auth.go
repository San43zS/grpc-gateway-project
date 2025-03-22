package middlewares

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-gateway-project/internal/models"
	"strings"
)

func (m *Middlewares) auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Missing authorization token")
		}
		authData := strings.Split(authHeader[0], " ")
		if len(authData) < 2 {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization token")
		}

		user, err := m.storage.User().GetUserByToken(ctx, authData[1])
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, models.UserCtxKey, user)

		return handler(ctx, req)
	}
}
