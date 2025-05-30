package grpc

import (
	"bmt_product_service/db/sqlc"
	"context"
	"fmt"
	rpc_product "product"
	"time"
)

type ProductRPCServer struct {
	SqlStore sqlc.Queries
	rpc_product.UnimplementedProductServer
}

func (p *ProductRPCServer) GetFilmDuration(ctx context.Context, arg *rpc_product.GetFilmDurationReq) (*rpc_product.GetFilmDurationRes, error) {
	filmDuration, err := p.SqlStore.GetFilmDuration(ctx, arg.FilmId)
	if err != nil {
		return nil, fmt.Errorf("failed to get film duration: %w", err)
	}

	duration := time.Duration(filmDuration.Microseconds) * time.Microsecond
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	s := int(duration.Seconds()) % 60

	return &rpc_product.GetFilmDurationRes{
		FilmDuration: fmt.Sprintf("%02dh%02dm%02ds", h, m, s),
	}, nil
}

func NewProductRPCServer() rpc_product.ProductServer {
	return &ProductRPCServer{}
}
