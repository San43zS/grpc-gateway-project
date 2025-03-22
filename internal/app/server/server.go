package server

import (
	"context"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"grpc-gateway-project/internal/app/server/grpcGwServer"
	"grpc-gateway-project/internal/app/server/grpcServer"
	"grpc-gateway-project/internal/config"
	"grpc-gateway-project/internal/handlers/grpcHandlers"
	"grpc-gateway-project/internal/scenarios"
	"grpc-gateway-project/internal/storage"
	"grpc-gateway-project/internal/storage/mongoDb"
	"log"
)

type server struct {
	config   *config.Config
	handlers *grpcHandlers.Handlers
	storage  storage.Storage
	ctx      context.Context
}

func New() *server {
	ctx := context.Background()

	storage, err := mongoDb.New(ctx)
	if err != nil {
		log.Fatalln("error in connection mongo : %w", err)
	}
	return &server{
		config:   config.New(),
		handlers: grpcHandlers.New(scenarios.New(storage)),
		storage:  storage,
		ctx:      ctx,
	}
}

func (s *server) Start() error {
	gr, ctx := errgroup.WithContext(s.ctx)

	gr.Go(func() error {
		return grpcGwServer.Start(ctx, s.config)
	})
	gr.Go(func() error {
		return grpcServer.New(s.handlers, s.config, s.storage).Start(ctx)
	})
	if err := gr.Wait(); err != nil {
		slog.Error(err.Error())
	}
	return nil
}
