package rpc

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

func (p *ProductRPCServer) GetPriceOfFAB(ctx context.Context, arg *rpc_product.GetPriceOfFABReq) (*rpc_product.GetPriceOfFABRes, error) {
	isFABExist, err := p.SqlStore.IsFilmExist(ctx, arg.FABId)
	if err != nil {
		return nil, fmt.Errorf("failed to check if fab exists or not: %w", err)
	}
	if !isFABExist {
		return nil, fmt.Errorf("fab with %d doesn't exist", arg.FABId)
	}

	price, err := p.SqlStore.GetPriceOfFABById(ctx, arg.FABId)
	if err != nil {
		return nil, fmt.Errorf("failed to get price of fab: %w", err)
	}

	return &rpc_product.GetPriceOfFABRes{
		Price: price,
	}, nil
}

func (p *ProductRPCServer) GetFilmDuration(ctx context.Context, arg *rpc_product.GetFilmDurationReq) (*rpc_product.GetFilmDurationRes, error) {
	isFilmExist, err := p.SqlStore.IsFilmExist(ctx, arg.FilmId)
	if err != nil {
		return nil, fmt.Errorf("failed to check if movie exists or not: %w", err)
	}
	if !isFilmExist {
		return nil, fmt.Errorf("film with %d doesn't exist", arg.FilmId)
	}

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

func NewProductRPCServer(
	sqlStore sqlc.Queries) rpc_product.ProductServer {
	return &ProductRPCServer{
		SqlStore: sqlStore,
	}
}
