package sqlc

import (
	"bmt_product_service/dto/request"
	"bmt_product_service/utils/convertors"
	"net/http"

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

// InsertFilmTran implements IStore.
func (s *SqlStore) InsertFilmTran(ctx context.Context, arg request.AddFilmReq) (int32, error) {
	var filmId int32 = -1

	err := s.execTran(ctx, func(q *Queries) error {
		interval, err := convertors.ParseDurationToPGInterval(arg.FilmInformation.Duration)
		if err != nil {
			return err
		}

		releaseDate, err := convertors.ConvertDateStringToTime(arg.FilmInformation.ReleaseDate)
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
			return fmt.Errorf("failed to insert film: %w", err)
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
			return fmt.Errorf("failed to insert film change: %w", err)
		}

		for _, genre := range arg.FilmInformation.Genres {
			var tmpGenre NullGenres
			err := tmpGenre.Scan(genre)
			if err != nil {
				return fmt.Errorf("failed to scan role: %w", err)
			}

			err = q.insertFilmGenre(ctx, insertFilmGenreParams{
				FilmID: pgtype.Int4{Int32: filmId, Valid: true},
				Genre: NullGenres{
					Genres: tmpGenre.Genres,
					Valid:  true,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to insert genre %s: %w", genre, err)
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
			return fmt.Errorf("failed to insert other film information: %w", err)
		}

		// payloadBytes, err := json.Marshal(gin.H{
		// 	"film_id":  filmId,
		// 	"duration": arg.FilmInformation.Duration,
		// })
		// if err != nil {
		// 	return fmt.Errorf("failed to marshal payload: %w", err)
		// }

		// err = q.CreateOutbox(ctx,
		// 	CreateOutboxParams{
		// 		AggregatedType: "PRODUCT_FILM_ID",
		// 		AggregatedID:   filmId,
		// 		EventType:      global.FILM_CREATED,
		// 		Payload:        payloadBytes,
		// 	})
		// if err != nil {
		// 	return fmt.Errorf("failed to create out box: %w", err)
		// }

		return nil
	})

	if err != nil {
		// If the transaction failed, return the error
		return -1, fmt.Errorf("transaction failed: %w", err)
	}

	return filmId, nil
}

// UpdateFilmTran implements IStore.
func (s *SqlStore) UpdateFilmTran(ctx context.Context, arg request.UpdateFilmReq) error {
	err := s.execTran(ctx, func(q *Queries) error {
		interval, err := convertors.ParseDurationToPGInterval(arg.FilmInformation.Duration)
		if err != nil {
			return err
		}

		releaseDate, err := convertors.ConvertDateStringToTime(arg.FilmInformation.ReleaseDate)
		if err != nil {
			return err
		}

		err = s.deleteAllFilmGenresByFilmID(ctx,
			pgtype.Int4{
				Int32: arg.FilmId,
				Valid: true,
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
					Int32: arg.FilmId,
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
			ID:          arg.FilmId,
			Title:       arg.FilmInformation.Title,
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
			FilmID:    arg.FilmId,
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

// InsertFABTran implements IStore.
func (s *SqlStore) InsertFABTran(ctx context.Context, arg request.AddFABReq) (int32, int, error) {
	var statusCode int = http.StatusOK
	var fABId int32 = -1

	err := s.execTran(ctx, func(q *Queries) error {
		var fabType NullFabTypes

		err := fabType.Scan(arg.Type)
		if err != nil {
			statusCode = http.StatusBadRequest

			return fmt.Errorf("invalid fab type: %v", err)
		}

		fABIdFromQuery, err := q.InsertFAB(ctx,
			InsertFABParams{
				Name: arg.Name,
				Type: fabType.FabTypes,
				ImageUrl: pgtype.Text{
					String: "",
					Valid:  true,
				},
				Price: int32(arg.Price),
			})
		if err != nil {
			statusCode = http.StatusInternalServerError

			return fmt.Errorf("insert FAB failed: %v", err)
		}

		fABId = fABIdFromQuery

		// payloadBytes, err := json.Marshal(gin.H{
		// 	"fab_id": fABId,
		// 	"price":  arg.Price,
		// })
		// if err != nil {
		// 	return fmt.Errorf("failed to update film change: %v", err)
		// }

		// err = q.CreateOutbox(ctx,
		// 	CreateOutboxParams{
		// 		AggregatedType: "PRODUCT_FAB_ID",
		// 		AggregatedID:   fABId,
		// 		EventType:      global.FAB_CREATED,
		// 		Payload:        payloadBytes,
		// 	})
		// if err != nil {
		// 	return fmt.Errorf("failed to create out box: %w", err)
		// }

		return nil
	})

	if err != nil {
		return fABId, statusCode, err
	}

	return fABId, statusCode, nil
}

func NewStore(connPool *pgxpool.Pool) IStore {
	return &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
