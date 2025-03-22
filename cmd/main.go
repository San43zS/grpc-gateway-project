package main

import (
	"grpc-gateway-project/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
