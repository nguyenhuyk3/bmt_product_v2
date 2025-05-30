-- name: GetFilmByTitle :one
SELECT *
FROM films
WHERE title = $1;

-- name: GetAllFilms :many
SELECT 
    f.id, f.title, f.description, f.release_date, f.duration,
    ARRAY_AGG(DISTINCT fg.genre::text) AS genres,
    ofi.status, ofi.poster_url, ofi.trailer_url
FROM films AS f
LEFT JOIN other_film_informations AS ofi ON f.id = ofi.film_id
LEFT JOIN film_genres AS fg ON fg.film_id = f.id
GROUP BY 
    f.id, f.title, f.description, f.release_date, f.duration,
    ofi.status, ofi.poster_url, ofi.trailer_url
ORDER BY f.release_date DESC;

-- name: GetPosterUrlByFilmId :one
SELECT poster_url
FROM other_film_informations
WHERE film_id = $1;

-- name: GetTrailerUrlByFilmId :one
SELECT trailer_url
FROM other_film_informations
WHERE film_id = $1;

-- name: IsFilmExist :one
SELECT EXISTS (
    SELECT 1 FROM films WHERE id = $1
) AS EXISTS;

-- name: GetFilmById :one
SELECT 
    f.id,
    f.title,
    f.description,
    f.release_date,
    f.duration,

    ofi.status,
    ofi.poster_url,
    ofi.trailer_url,

    ARRAY_AGG(DISTINCT fg.genre::text) AS genres

FROM films AS f
LEFT JOIN other_film_informations ofi ON ofi.film_id = f.id
LEFT JOIN film_genres fg ON fg.film_id = f.id

WHERE f.id = $1

GROUP BY 
    f.id,
    ofi.status,
    ofi.poster_url,
    ofi.trailer_url;

-- name: GetFilmDuration :one
SELECT duration
FROM films
WHERE id = $1;

