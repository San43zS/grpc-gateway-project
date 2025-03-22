package app

import (
	"github.com/joho/godotenv"
	"grpc-gateway-project/internal/app/server"
)

func Start() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := server.New().Start(); err != nil {
		return err
	}

	return nil
}
