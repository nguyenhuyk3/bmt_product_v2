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

type filmService struct {
	UploadService services.IUpload
	SqlStore      sqlc.IStore
	RedisClient   services.IRedis
}

func NewFilmService(
	uploadService services.IUpload,
	sqlStore sqlc.IStore,
	redisClient services.IRedis) services.IFilm {
	return &filmService{
		UploadService: uploadService,
		SqlStore:      sqlStore,
		RedisClient:   redisClient,
	}
}

// AddFilm implements services.IFilm.
func (f *filmService) AddFilm(ctx context.Context, arg request.AddFilmReq) (int, error) {
	filmId, err := f.SqlStore.InsertFilmTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	go func() {
		f.RedisClient.Save(fmt.Sprintf("%s%s", global.FILM, filmId), filmId, 60*24*30)
	}()

	go func() {
		err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
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
		err := f.UploadService.UploadFilmVideoToS3(request.UploadVideoReq{
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
func (f *filmService) GetAllFilms(ctx context.Context) (int, interface{}, error) {
	var films []sqlc.GetAllFilmsRow

	err := f.RedisClient.Get(global.GET_ALL_FILMS_WITH_ADMIN_ROLE, &films)
	if err != nil {
		if err.Error() == fmt.Sprintf("key %s does not exist", global.GET_ALL_FILMS_WITH_ADMIN_ROLE) {
			films, err = f.SqlStore.GetAllFilms(ctx)
			if err != nil {
				return http.StatusInternalServerError, nil, err
			}

			savingErr := f.RedisClient.Save(global.GET_ALL_FILMS_WITH_ADMIN_ROLE, &films, 60*24*10)
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
func (f *filmService) UpdateFilm(ctx context.Context, arg request.UpdateFilmReq) (int, error) {
	filmId, err := strconv.Atoi(arg.FilmId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("invalid film id (%s): %v", arg.FilmId, err)
	}

	isExist, err := f.SqlStore.IsFilmExist(ctx, int32(filmId))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("an error occur when querying database: %v", err)
	}
	if !isExist {
		return http.StatusNotFound, fmt.Errorf("film doesn't exist")
	}

	if arg.OtherFilmInformation.PosterFile != nil {
		go func() {
			err := f.UploadService.UploadProductImageToS3(request.UploadImageReq{
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

				objectKey, err := f.SqlStore.GetPosterUrlByFilmId(context.Background(), int32(filmId))
				if err != nil {
					log.Printf("failed to get poster URL: %d %v\n", filmId, err)
					return
				}

				if objectKey.String == "" {
					log.Println("no poster URL to delete (film poster)")
					return
				}

				err = f.UploadService.DeleteObject(objectKey.String)
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
			err := f.UploadService.UploadFilmVideoToS3(request.UploadVideoReq{
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

				objectKey, err := f.SqlStore.GetTrailerUrlByFilmId(context.Background(), int32(filmId))
				if err != nil {
					log.Printf("failed to get trailer URL: %d %v\n", filmId, err)
					return
				}

				if objectKey.String == "" {
					log.Println("no trailer URL to delete (film video)")
					return
				}

				err = f.UploadService.DeleteObject(objectKey.String)
				if err != nil {
					log.Printf("failed to delete trailer from S3: %v\n", err)
					return
				}

				fmt.Println("trailer deleted successfully")
			}
		}()
	}

	err = f.SqlStore.UpdateFilmTran(ctx, arg)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
