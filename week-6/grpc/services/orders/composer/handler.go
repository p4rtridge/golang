package composer

import (
	"kitchen/services/common/genproto/orders"
	"kitchen/services/orders/handler"
	"kitchen/services/orders/service"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func ServeRPC(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	rpcServer := grpc.NewServer()

	service := service.NewOrderService()

	handler := handler.NewOrdersRPCHandler(service)

	orders.RegisterOrderServiceServer(rpcServer, handler)

	log.Println("RPC server is listening on: ", addr)

	return rpcServer.Serve(lis)
}

func ServeHTTP(addr string) error {
	router := http.NewServeMux()

	service := service.NewOrderService()

	handler := handler.NewOrdersHTTPHandler(service)

	router.HandleFunc("POST /orders", handler.CreateOrder)

	log.Println("HTTP server is listening on: ", addr)

	return http.ListenAndServe(addr, router)
}
