package rpc

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/utils/convertors"
	"context"
	"fmt"
	rpc_product "product"
	"strings"
	"time"
)

type ProductRPCServer struct {
	SqlStore sqlc.Queries
	rpc_product.UnimplementedProductServer
}

func (p *ProductRPCServer) GetFilmCurrentlyShowing(ctx context.Context, arg *rpc_product.GetFilmCurrentlyShowingReq) (*rpc_product.GetFilmCurrentlyShowingRes, error) {
	film, err := p.SqlStore.GetFilmById(ctx, arg.FilmId)
	if err != nil {
		return nil, fmt.Errorf("failed to get film info for showing: %w", err)
	}

	genres, err := convertors.ConvertInterfaceToSlice(film.Genres)
	if err != nil {
		return nil, fmt.Errorf("failed to convert interface to slice: %w", err)
	}

	duration := time.Duration(film.Duration.Microseconds) * time.Microsecond
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60

	return &rpc_product.GetFilmCurrentlyShowingRes{
		FilmId:    film.ID,
		PosterUrl: film.PosterUrl.String,
		Genres:    strings.Join(genres, ", "),
		Duration:  fmt.Sprintf("%02dh%02dm", h, m),
	}, nil
}

func (p *ProductRPCServer) CheckFABExist(ctx context.Context, arg *rpc_product.CheckFABExistReq) (*rpc_product.CheckFABExistRes, error) {
	isFABExist, err := p.SqlStore.IsFilmExist(ctx, arg.FABId)
	if err != nil {
		return nil, fmt.Errorf("failed to check if fab exists or not: %w", err)
	}
	if !isFABExist {
		return nil, fmt.Errorf("fab with %d doesn't exist", arg.FABId)
	}

	return &rpc_product.CheckFABExistRes{ResponseMessage: "fab is exists"}, nil
}

func (p *ProductRPCServer) GetPriceOfFAB(ctx context.Context, arg *rpc_product.GetPriceOfFABReq) (*rpc_product.GetPriceOfFABRes, error) {
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
