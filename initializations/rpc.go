package initializations

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/global"
	"bmt_product_service/internal/rpc"
	"fmt"
	"log"
	"net"

	rpc_product "product"

	"google.golang.org/grpc"
)

func initRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", global.Config.Server.RPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	productRPCServer := rpc.NewProductRPCServer(*sqlc.New(global.Postgresql))
	grpcServer := grpc.NewServer()

	rpc_product.RegisterProductServer(grpcServer, productRPCServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
