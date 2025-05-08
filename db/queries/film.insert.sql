-- name: insertFilm :one
INSERT INTO "films" ("title", "description", "release_date", "duration")
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: insertFilmChange :exec
INSERT INTO "fillm_changes" ("film_id", "changed_by", "created_at", "updated_at")
VALUES ($1, $2, $3, $4);

-- name: insertFilmGenre :exec
INSERT INTO "film_genres" (film_id, genre)
VALUES ($1, $2)
ON CONFLICT (film_id, genre) DO NOTHING;

-- name: insertOtherFilmInformation :exec 
INSERT INTO "other_film_informations" ("film_id","status", "poster_url", "trailer_url")
VALUES ($1, $2, $3, $4);