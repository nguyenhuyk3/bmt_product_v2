package sqlc

import (
	"bmt_product_service/dto/request"
	"context"
)

type IStore interface {
	Querier
	InsertFilmTran(ctx context.Context, arg request.AddFilmReq) (string, error)
	UpdateFilmTran(ctx context.Context, arg request.UpdateFilmReq) error
}
