package middlewares

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"grpc-gateway-project/internal/storage"
)

type Middlewares struct {
	storage storage.Storage
}

func New(storage storage.Storage) grpc.UnaryServerInterceptor {
	m := &Middlewares{
		storage: storage,
	}
	return configureMainInterceptor(m.configureMethodInterceptors("api"))
}

func (m *Middlewares) configureMethodInterceptors(packageName string) map[string]grpc.UnaryServerInterceptor {
	base := "/" + packageName + "."
	return map[string]grpc.UnaryServerInterceptor{
		base + "UserService/SubscribeUser":   grpc_middleware.ChainUnaryServer(m.auth()),
		base + "UserService/UnsubscribeUser": grpc_middleware.ChainUnaryServer(m.auth()),
		base + "UserService/DeleteUser":      grpc_middleware.ChainUnaryServer(m.auth()),
	}
}

func configureMainInterceptor(methodInterceptors map[string]grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if interceptor, ok := methodInterceptors[info.FullMethod]; ok {
			return interceptor(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}
