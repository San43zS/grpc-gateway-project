package grpcGwServer

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-gateway-project/internal/config"
	"grpc-gateway-project/proto/api/generate/desc"
	"log"
	"net/http"
	"strconv"
)

func Start(ctx context.Context, config *config.Config) error {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:"+strconv.Itoa(config.GrpcPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	gwMux := runtime.NewServeMux()

	err = desc.RegisterUserServiceHandler(context.Background(), gwMux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway : ", err)
	}
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GateWayPort),
		Handler: setGrpcMetadata(gwMux),
	}
	log.Println("serving grpc-gateway on http://localhost:8000")

	gwServer.RegisterOnShutdown(func() {
		err = conn.Close()
		slog.Error(err.Error())
	})
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		errCh <- gwServer.ListenAndServe()
	}()
	select {
	case err = <-errCh:
	case <-ctx.Done():
		err = gwServer.Shutdown(ctx)
		slog.Error(err.Error())
	}
	return err
}
