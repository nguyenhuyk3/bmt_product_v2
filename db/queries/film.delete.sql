-- name: deleteAllFilmGenresByFilmID :exec
DELETE FROM film_genres
WHERE film_id = $1;
