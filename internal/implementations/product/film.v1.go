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
)

type productService struct {
	UploadService services.IUpload
	SqlStore      sqlc.IStore
	RedisClient   services.IRedis
}

func NewProductService(
	uploadService services.IUpload,
	sqlStore sqlc.IStore,
	redisClient services.IRedis) services.IFilm {
	return &productService{
		UploadService: uploadService,
		SqlStore:      sqlStore,
		RedisClient:   redisClient,
	}
}

// AddFilm implements services.IFilm.
func (p *productService) AddFilm(ctx context.Context, arg request.AddFilmReq) (int, error) {
	filmId, err := p.SqlStore.InsertFilmTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	go func() {
		p.RedisClient.Save(fmt.Sprintf("%s%s", global.FILM, filmId), filmId, 60*24*30)
	}()

	go func() {
		err := p.UploadService.UploadProductImageToS3(request.UploadImageReq{
			ProductId: filmId,
			Image:     arg.OtherFilmInformation.PosterFile,
		}, global.FILM_TYPE)
		if err != nil {
			log.Printf("an error occurr when updating film poster to S3 (add): %v", err)
		} else {
			log.Printf("upload film poster to S3 successfully")
		}
	}()

	go func() {
		err := p.UploadService.UploadFilmVideoToS3(request.UploadVideoReq{
			ProductId: filmId,
			Video:     arg.OtherFilmInformation.TrailerFile,
		})
		if err != nil {
			log.Printf("an error occurr when updating film trailer to S3 (add): %v", err)
		} else {
			log.Printf("upload film trailer to S3 successfully")
		}
	}()

	return http.StatusOK, nil
}

// GetAllFilms implements services.IFilm.
func (p *productService) GetAllFilms(ctx context.Context) (int, interface{}, error) {
	var films []sqlc.GetAllFilmsRow

	err := p.RedisClient.Get(global.GET_ALL_FILMS_WITH_ADMIN_ROLE, &films)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", global.GET_ALL_FILMS_WITH_ADMIN_ROLE) {
			films, err = p.SqlStore.GetAllFilms(ctx)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}

			savingErr := p.RedisClient.Save(global.GET_ALL_FILMS_WITH_ADMIN_ROLE, &films, 60*24*10)
			if savingErr != nil {
				return http.StatusOK, nil, fmt.Errorf("warning: failed to save to Redis: %v", savingErr)
			}

			return http.StatusOK, films, nil
		}

		return http.StatusInternalServerError, nil, fmt.Errorf("getting value occur an error: %v", err)
	}

	return http.StatusOK, films, nil
}

// UpdateFilm implements services.IFilm.
func (p *productService) UpdateFilm(ctx context.Context, arg request.UpdateFilmReq) (int, error) {
	if arg.OtherFilmInformation.PosterFile != nil {
		go func() {
			err := p.UploadService.UploadProductImageToS3(request.UploadImageReq{
				ProductId: arg.FilmId,
				Image:     arg.OtherFilmInformation.PosterFile,
			}, global.FILM_TYPE)
			if err != nil {
				log.Printf("an error occurr when updating film poster to S3 (update): %v", err)
			} else {
				filmId, err := strconv.Atoi(arg.FilmId)
				if err != nil {
					log.Printf("invalid film ID: %v\n", err)
					return
				}

				objectKey, err := p.SqlStore.GetPosterUrlByFilmId(ctx, int32(filmId))
				if err != nil {
					log.Printf("failed to get poster URL: %v\n", err)
					return
				}

				if objectKey.String == "" {
					log.Println("no poster URL to delete")
					return
				}

				err = p.UploadService.DeleteObject(objectKey.String)
				if err != nil {
					log.Printf("failed to delete poster from S3: %v\n", err)
					return
				}

				fmt.Println("poster deleted successfully")
			}
		}()

	}

	if arg.OtherFilmInformation.TrailerFile != nil {
		go func() {
			err := p.UploadService.UploadFilmVideoToS3(request.UploadVideoReq{
				ProductId: arg.FilmId,
				Video:     arg.OtherFilmInformation.TrailerFile,
			})
			if err != nil {
				log.Printf("an error occurr when updating film trailer to S3 (update): %v", err)
			} else {
				filmId, err := strconv.Atoi(arg.FilmId)
				if err != nil {
					log.Printf("invalid film ID: %v\n", err)
					return
				}

				objectKey, err := p.SqlStore.GetTrailerUrlByFilmId(ctx, int32(filmId))
				if err != nil {
					log.Printf("failed to get trailer URL: %v\n", err)
					return
				}

				if objectKey.String == "" {
					log.Println("no trailer URL to delete")
					return
				}

				err = p.UploadService.DeleteObject(objectKey.String)
				if err != nil {
					log.Printf("failed to delete trailer from S3: %v\n", err)
					return
				}

				fmt.Println("trailer deleted successfully")
			}
		}()
	}

	err := p.SqlStore.UpdateFilmTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
