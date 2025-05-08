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
