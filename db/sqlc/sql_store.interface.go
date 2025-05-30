package sqlc

import (
	"bmt_product_service/dto/request"
	"context"
)

type IStore interface {
	Querier
	InsertFilmTran(ctx context.Context, arg request.AddFilmReq) (int32, error)
	UpdateFilmTran(ctx context.Context, arg request.UpdateFilmReq) error
	InsertFABTran(ctx context.Context, arg request.AddFABReq) (int32, int, error)
}
