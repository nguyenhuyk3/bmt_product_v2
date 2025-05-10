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
			log.Printf("an error occur when uploading image to S3 (fab): %v", err)
		} else {
			log.Println("upload image to S3 successfully (fab)")
		}
	}()

	return http.StatusOK, nil
}

// DeleteFAB implements services.IFoodAndBeverage.
func (f *fABService) DeleteFAB(ctx context.Context, fABId int32) (int, error) {
	err := f.SqlStore.ToggleFABDelete(ctx, fABId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("an error occur when toggling fab: %v", err)
	}
	return http.StatusOK, nil
}

// GetAllFAB implements services.IFoodAndBeverage.
func (f *fABService) GetAllFAB(ctx context.Context) (interface{}, int, error) {
	panic("unimplemented")
}

// UpdateFAB implements services.IFoodAndBeverage.
func (f *fABService) UpdateFAB(ctx context.Context, arg request.UpdateFABReq) (int, error) {
	if arg.Image != nil {
		go func() {
			err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
				ProductId: strconv.Itoa(int(arg.FABId)),
				Image:     arg.Image,
			}, global.FAB_TYPE)
			if err != nil {
				log.Printf("an error occur when upading image (fab): %v", err)
			} else {
				log.Println("upload image to S3 successfully (fab)")
			}
		}()
	}
	var fabType sqlc.NullFabTypes
	err := fabType.Scan(arg.Type)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid fab type: %v", err)
	}

	err = f.SqlStore.UpdateFAB(ctx, sqlc.UpdateFABParams{
		ID:   arg.FABId,
		Name: arg.Name,
		Type: fabType.FabTypes,
		ImageUrl: pgtype.Text{
			String: "",
			Valid:  true,
		},
		Price: int32(arg.Price),
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("an error occur when updating fab product: %v", err)
	}

	return http.StatusOK, nil
}

func NewFABService(
	uploadService services.IUpload,
	sqlStore sqlc.IStore,
) services.IFoodAndBeverage {
	return &fABService{
		UploadService: uploadService,
		SqlStore:      sqlStore,
	}
}
