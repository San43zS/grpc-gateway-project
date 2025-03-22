DOCKER_COMPOSE = docker-compose

#grpc compile config
in_grpc_dir=proto/api
out_grpc_dir=proto/api
imports_grpc_dir=./proto

--grpc-gateway_out=logtostderr=true:. \

test:
	go test ./internal/scenarios -coverprofile=coverage.out /

compile_proto:
	protoc -I$(imports_grpc_dir) --proto_path=$(in_grpc_dir) \
	--go_out=$(out_grpc_dir) \
	--go-grpc_out=$(out_grpc_dir) \
	--grpc-gateway_out=allow_delete_body=true:$(out_grpc_dir) \
	$(in_grpc_dir)/api.proto

docker-compose-up:
	$(DOCKER_COMPOSE) up --build

# Default target
default: test