package initializations

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/global"
	"bmt_product_service/internal/rpc"
	"log"
	"net"

	rpc_product "product"

	"google.golang.org/grpc"
)

func initRPC() {
	lis, err := net.Listen("tcp", ":50033")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	productRPCServer := rpc.NewProductRPCServer(*sqlc.New(global.Postgresql))

	grpcServer := grpc.NewServer()
	rpc_product.RegisterProductServer(grpcServer, productRPCServer)

	log.Println("Product Service listening on :50033")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
