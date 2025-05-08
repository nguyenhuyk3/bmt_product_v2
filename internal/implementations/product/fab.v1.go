package product

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/dto/request"
	"bmt_product_service/global"
	"bmt_product_service/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

type fABService struct {
	UploadService services.IUpload
	SqlStore      sqlc.IStore
	RedisClient   services.IRedis
}

// AddFAB implements services.IFoodAndBeverage.
func (f *fABService) AddFAB(ctx context.Context, arg request.AddFABReq) (int, error) {
	var fabType sqlc.NullFabTypes
	err := fabType.Scan(arg.Type)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid fab type: %v", err)
	}

	fABId, err := f.SqlStore.InsertFAB(ctx,
		sqlc.InsertFABParams{
			Name: arg.Name,
			Type: fabType.FabTypes,
			ImageUrl: pgtype.Text{
				String: "",
				Valid:  true,
			},
			Price: int32(arg.Price),
		})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("insert FAB failed: %v", err)
	}

	go func() {
		err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
			ProductId: strconv.Itoa(int(fABId)),
			Image:     arg.Image,
		}, global.FAB_TYPE)
		if err != nil {
			log.Printf("an error occurr when uploading image to S3 (fab): %v", err)
		} else {
			log.Println("upload image to S3 successfully (fab)")
		}
	}()

	return http.StatusOK, nil
}

// DeleteFAB implements services.IFoodAndBeverage.
func (f *fABService) DeleteFAB(ctx context.Context, fABId int32) (int, error) {

	return http.StatusOK, nil
}

// GetAllFAB implements services.IFoodAndBeverage.
func (f *fABService) GetAllFAB(ctx context.Context) (interface{}, int, error) {
	panic("unimplemented")
}

// UpdateFAB implements services.IFoodAndBeverage.
func (f *fABService) UpdateFAB(ctx context.Context, arg request.UpdateFABReq) (int, error) {
	panic("unimplemented")
}

func NewFABService(
	uploadService services.IUpload,
	sqlStore sqlc.IStore,
	redisClient services.IRedis) services.IFoodAndBeverage {
	return &fABService{
		UploadService: uploadService,
		SqlStore:      sqlStore,
		RedisClient:   redisClient,
	}
}
