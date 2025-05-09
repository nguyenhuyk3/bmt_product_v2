// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: film.insert.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertFilm = `-- name: insertFilm :one
INSERT INTO "films" ("title", "description", "release_date", "duration")
VALUES ($1, $2, $3, $4)
RETURNING id
`

type insertFilmParams struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ReleaseDate pgtype.Date     `json:"release_date"`
	Duration    pgtype.Interval `json:"duration"`
}

func (q *Queries) insertFilm(ctx context.Context, arg insertFilmParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertFilm,
		arg.Title,
		arg.Description,
		arg.ReleaseDate,
		arg.Duration,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const insertFilmChange = `-- name: insertFilmChange :exec
INSERT INTO "fillm_changes" ("film_id", "changed_by", "created_at", "updated_at")
VALUES ($1, $2, $3, $4)
`

type insertFilmChangeParams struct {
	FilmID    int32            `json:"film_id"`
	ChangedBy string           `json:"changed_by"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) insertFilmChange(ctx context.Context, arg insertFilmChangeParams) error {
	_, err := q.db.Exec(ctx, insertFilmChange,
		arg.FilmID,
		arg.ChangedBy,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const insertFilmGenre = `-- name: insertFilmGenre :exec
INSERT INTO "film_genres" (film_id, genre)
VALUES ($1, $2)
ON CONFLICT (film_id, genre) DO NOTHING
`

type insertFilmGenreParams struct {
	FilmID pgtype.Int4 `json:"film_id"`
	Genre  NullGenres  `json:"genre"`
}

func (q *Queries) insertFilmGenre(ctx context.Context, arg insertFilmGenreParams) error {
	_, err := q.db.Exec(ctx, insertFilmGenre, arg.FilmID, arg.Genre)
	return err
}

const insertOtherFilmInformation = `-- name: insertOtherFilmInformation :exec
INSERT INTO "other_film_informations" ("film_id","status", "poster_url", "trailer_url")
VALUES ($1, $2, $3, $4)
`

type insertOtherFilmInformationParams struct {
	FilmID     int32        `json:"film_id"`
	Status     NullStatuses `json:"status"`
	PosterUrl  pgtype.Text  `json:"poster_url"`
	TrailerUrl pgtype.Text  `json:"trailer_url"`
}

func (q *Queries) insertOtherFilmInformation(ctx context.Context, arg insertOtherFilmInformationParams) error {
	_, err := q.db.Exec(ctx, insertOtherFilmInformation,
		arg.FilmID,
		arg.Status,
		arg.PosterUrl,
		arg.TrailerUrl,
	)
	return err
}
