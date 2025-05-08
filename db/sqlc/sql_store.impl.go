package sqlc

import (
	"bmt_product_service/dto/request"
	"bmt_product_service/utils/convertors"
	"strconv"

	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func (s *SqlStore) execTran(ctx context.Context, fn func(*Queries) error) error {
	// Start transaction
	tran, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tran)
	// fn performs a series of operations down the db
	err = fn(q)
	if err != nil {
		// If an error occurs, rollback the transaction
		if rbErr := tran.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tran err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tran.Commit(ctx)
}

func parseDurationToPGInterval(durationStr string) (pgtype.Interval, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return pgtype.Interval{}, fmt.Errorf("invalid duration format: %v", err)
	}

	return pgtype.Interval{
		Microseconds: duration.Microseconds(),
		Valid:        true,
	}, nil
}

// InsertFilmTran implements IStore.
func (s *SqlStore) InsertFilmTran(ctx context.Context, arg request.AddFilmReq) (string, error) {
	var filmId int32 = -1

	err := s.execTran(ctx, func(q *Queries) error {
		interval, err := parseDurationToPGInterval(arg.FilmInformation.Duration)
		if err != nil {
			return err
		}

		releaseDate, err := convertors.GetReleaseDateAsTime(arg.FilmInformation.ReleaseDate)
		if err != nil {
			return err
		}

		filmId, err = q.insertFilm(ctx, insertFilmParams{
			Title:       arg.FilmInformation.Title,
			Description: arg.FilmInformation.Description,
			ReleaseDate: pgtype.Date{
				Time:  releaseDate,
				Valid: true,
			},
			Duration: interval,
		})
		if err != nil {
			return fmt.Errorf("failed to insert film: %v", err)
		}

		err = q.insertFilmChange(ctx, insertFilmChangeParams{
			FilmID:    filmId,
			ChangedBy: arg.ChangedBy,
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to insert film change: %v", err)
		}

		for _, genre := range arg.FilmInformation.Genres {
			var tmpGenre NullGenres
			err := tmpGenre.Scan(genre)
			if err != nil {
				return fmt.Errorf("failed to scan role: %v", err)
			}

			err = q.insertFilmGenre(ctx, insertFilmGenreParams{
				FilmID: pgtype.Int4{Int32: filmId, Valid: true},
				Genre: NullGenres{
					Genres: tmpGenre.Genres,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to insert genre %s: %v", genre, err)
			}
		}

		// var filmStatus NullStatuses
		// if err = filmStatus.Scan(arg.OtherFilmInformation.Status); err != nil {
		// 	return fmt.Errorf("failed to scan status: %v", err)
		// }

		err = q.insertOtherFilmInformation(ctx, insertOtherFilmInformationParams{
			FilmID: filmId,
			Status: NullStatuses{
				Statuses: StatusesPending,
				Valid:    true,
			},
			PosterUrl: pgtype.Text{
				String: "",
				Valid:  true,
			},
			TrailerUrl: pgtype.Text{
				String: "",
				Valid:  true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to insert other film information: %v", err)
		}

		return nil
	})

	if err != nil {
		// If the transaction failed, return the error
		return "", fmt.Errorf("transaction failed: %v", err)
	}

	return strconv.Itoa(int(filmId)), nil
}

// UpdateFilmTran implements IStore.
func (s *SqlStore) UpdateFilmTran(ctx context.Context, arg request.UpdateFilmReq) error {
	err := s.execTran(ctx, func(q *Queries) error {
		interval, err := parseDurationToPGInterval(arg.FilmInformation.Duration)
		if err != nil {
			return err
		}

		releaseDate, err := convertors.GetReleaseDateAsTime(arg.FilmInformation.ReleaseDate)
		if err != nil {
			return err
		}

		filmId, err := strconv.Atoi(arg.FilmId)
		if err != nil {
			return fmt.Errorf("invalid filmId '%s': %v", arg.FilmId, err)
		}

		err = s.deleteAllFilmGenresByFilmID(ctx, pgtype.Int4{
			Int32: int32(filmId), Valid: true,
		})
		if err != nil {
			return fmt.Errorf("failed to delete existing genres: %v", err)
		}

		for _, genre := range arg.FilmInformation.Genres {
			var tmpGenre NullGenres
			err := tmpGenre.Scan(genre)
			if err != nil {
				return fmt.Errorf("failed to scan role: %v", err)
			}

			err = q.insertFilmGenre(ctx, insertFilmGenreParams{
				FilmID: pgtype.Int4{
					Int32: int32(filmId),
					Valid: true,
				},
				Genre: NullGenres{
					Genres: tmpGenre.Genres,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to insert genre %s: %v", genre, err)
			}
		}

		err = s.updateFilm(ctx, updateFilmParams{
			ID:          int32(filmId),
			Description: arg.FilmInformation.Description,
			ReleaseDate: pgtype.Date{
				Time:  releaseDate,
				Valid: true,
			},
			Duration: interval,
		})
		if err != nil {
			return fmt.Errorf("failed to update film: %v", err)
		}

		err = s.updateFilmChange(ctx, updateFilmChangeParams{
			FilmID:    int32(filmId),
			ChangedBy: arg.ChangedBy,
			UpdatedAt: pgtype.Timestamp{
				Time: time.Now(), Valid: true,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update film change: %v", err)
		}

		return nil
	})

	if err != nil {
		// If the transaction failed, return the error
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}

func NewStore(connPool *pgxpool.Pool) IStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
