package composer

import (
	"kitchen/services/kitchen/repository/rpc"
	"kitchen/services/kitchen/transport"
	"log"
	"net/http"
)

func ServeHTTP(addr string) error {
	router := http.NewServeMux()

	repo := rpc.NewGRPCClient("0.0.0.0:50051")

	api := transport.NewHTTPTransport(repo)

	router.HandleFunc("/", api.ServeHomepage)

	log.Println("Starting server on", addr)
	return http.ListenAndServe(addr, router)
}
