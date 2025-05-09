// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	GetAllFilms(ctx context.Context) ([]GetAllFilmsRow, error)
	GetFABById(ctx context.Context, id int32) (FoodAndBeverage, error)
	GetFABImageURLByID(ctx context.Context, id int32) (pgtype.Text, error)
	GetFilmByTitle(ctx context.Context, title string) (Films, error)
	GetPosterUrlByFilmId(ctx context.Context, filmID int32) (pgtype.Text, error)
	GetTrailerUrlByFilmId(ctx context.Context, filmID int32) (pgtype.Text, error)
	InsertFAB(ctx context.Context, arg InsertFABParams) (int32, error)
	ListFAB(ctx context.Context) ([]FoodAndBeverage, error)
	ToggleFABDelete(ctx context.Context, id int32) error
	UpdateFAB(ctx context.Context, arg UpdateFABParams) error
	UpdateFABImageURL(ctx context.Context, arg UpdateFABImageURLParams) error
	UpdatePosterUrlAndCheckStatus(ctx context.Context, arg UpdatePosterUrlAndCheckStatusParams) error
	UpdateVideoUrlAndCheckStatus(ctx context.Context, arg UpdateVideoUrlAndCheckStatusParams) error
	deleteAllFilmGenresByFilmID(ctx context.Context, filmID pgtype.Int4) error
	insertFilm(ctx context.Context, arg insertFilmParams) (int32, error)
	insertFilmChange(ctx context.Context, arg insertFilmChangeParams) error
	insertFilmGenre(ctx context.Context, arg insertFilmGenreParams) error
	insertOtherFilmInformation(ctx context.Context, arg insertOtherFilmInformationParams) error
	updateFilm(ctx context.Context, arg updateFilmParams) error
	updateFilmChange(ctx context.Context, arg updateFilmChangeParams) error
	updateFilmInformation(ctx context.Context, arg updateFilmInformationParams) error
}

var _ Querier = (*Queries)(nil)
