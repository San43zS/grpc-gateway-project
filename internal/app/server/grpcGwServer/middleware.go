package grpcGwServer

import (
	"net/http"
)

func setGrpcMetadata(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		handler.ServeHTTP(w, r)
	})
}
