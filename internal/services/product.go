package services

import (
	"bmt_product_service/dto/request"
	"context"
)

type IFilm interface {
	AddFilm(ctx context.Context, arg request.AddFilmReq) (int, error)
	UpdateFilm(ctx context.Context, arg request.UpdateFilmReq) (int, error)
	GetAllFilms(ctx context.Context) (int, interface{}, error)
}

type IFoodAndBeverage interface {
	AddFAB(ctx context.Context, arg request.AddFABReq) (int, error)
	UpdateFAB(ctx context.Context, arg request.UpdateFABReq) (int, error)
	DeleteFAB(ctx context.Context, fABId int32) (int, error)
	GetAllFAB(ctx context.Context) (interface{}, int, error)
}
