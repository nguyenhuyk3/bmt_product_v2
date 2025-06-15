package product

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/dto/request"
	"bmt_product_service/dto/response"
	"bmt_product_service/global"
	"bmt_product_service/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
)

type fABService struct {
	UploadService services.IUpload
	SqlStore      sqlc.IStore
	RedisClient   services.IRedis
}

// AddFAB implements services.IFoodAndBeverage.
func (f *fABService) AddFAB(ctx context.Context, arg request.AddFABReq) (int, error) {
	fABId, status, err := f.SqlStore.InsertFABTran(ctx, arg)
	if err != nil {
		return status, err
	}

	go func() {
		err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
			ProductId: fABId,
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

// GetAllFABs implements services.IFoodAndBeverage.
func (f *fABService) GetAllFABs(ctx context.Context) (interface{}, int, error) {
	var fabs []response.FABItem

	err := f.RedisClient.Get(global.ALL_FABS, &fabs)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", global.ALL_FABS) {
			queryedFABs, err := f.SqlStore.GetAllFABs(ctx)
			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("an error occur when querying fabs: %w", err)
			}

			for _, fAB := range queryedFABs {
				fabs = append(fabs, response.FABItem{
					ID:       fAB.ID,
					Name:     fAB.Name,
					Type:     string(fAB.Type),
					ImageUrl: fAB.ImageUrl.String,
					Price:    fAB.Price,
				})
			}

			err = f.RedisClient.Save(global.ALL_FABS, fabs, thirty_days)
			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("warning: failed to save to Redis: %w", err)
			}

			return fabs, http.StatusOK, nil
		}

		return nil, http.StatusInternalServerError, fmt.Errorf("getting value occur an error: %w", err)
	}

	return fabs, http.StatusOK, nil
}

// UpdateFAB implements services.IFoodAndBeverage.
func (f *fABService) UpdateFAB(ctx context.Context, arg request.UpdateFABReq) (int, error) {
	isExist, err := f.SqlStore.IsFABExist(ctx, arg.FABId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("an error occur when querying database: %v", err)
	}
	if !isExist {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	if arg.Image != nil {
		go func() {
			err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
				ProductId: arg.FABId,
				Image:     arg.Image,
			}, global.FAB_TYPE)
			if err != nil {
				log.Printf("an error occur when upading image (fab): %v", err)
			} else {
				objectKey, err := f.SqlStore.GetFABImageURLByID(context.Background(), arg.FABId)
				if err != nil {
					log.Printf("failed to get image URL (fab): %d %v\n", arg.FABId, err)
					return
				}

				if objectKey.String == "" {
					log.Printf("no image URL to delete (fab) with id: %d", arg.FABId)
					return
				}

				err = f.UploadService.DeleteObject(objectKey.String)
				if err != nil {
					log.Printf("failed to delete image from S3 (fab): %v\n", err)
					return
				}

				log.Println("upload image to S3 successfully (fab)")
			}
		}()
	}
	var fabType sqlc.NullFabTypes
	err = fabType.Scan(arg.Type)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid fab type: %v", err)
	}

	err = f.SqlStore.UpdateFAB(ctx, sqlc.UpdateFABParams{
		ID:    arg.FABId,
		Name:  arg.Name,
		Type:  fabType.FabTypes,
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
	redisClient services.IRedis,
) services.IFoodAndBeverage {
	return &fABService{
		UploadService: uploadService,
		SqlStore:      sqlStore,
		RedisClient:   redisClient,
	}
}
