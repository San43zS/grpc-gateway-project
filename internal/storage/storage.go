package storage

import "grpc-gateway-project/internal/storage/repsInterfaces"

type Storage interface {
	User() repsInterfaces.User
}
