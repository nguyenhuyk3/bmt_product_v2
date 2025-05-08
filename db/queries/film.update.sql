-- name: UpdatePosterUrlAndCheckStatus :exec
UPDATE "other_film_informations"
SET poster_url = $2, 
    status = CASE 
        WHEN trailer_url IS NOT NULL
            AND LENGTH(trailer_url) > 0 
            AND LENGTH($2::text) > 0 THEN 'success' 
        ELSE status
    END
WHERE "film_id" = $1;

-- name: UpdateVideoUrlAndCheckStatus :exec
UPDATE "other_film_informations"
SET trailer_url = $2, 
    status = CASE 
        WHEN poster_url IS NOT NULL 
        AND LENGTH(poster_url) > 0
        AND LENGTH($2::text) > 0 THEN 'success' 
        ELSE status
    END
WHERE "film_id" = $1;

-- name: updateFilm :exec
UPDATE films
SET 
    title = $2,
    description = $3,
    release_date = $4,
    duration = $5
WHERE id = $1;

-- name: updateFilmChange :exec
UPDATE fillm_changes
SET 
    changed_by = $2,
    updated_at = $3
WHERE film_id = $1;

-- name: updateFilmInformation :exec
UPDATE other_film_informations
SET 
    status = $2,
    poster_url = $3,
    trailer_url = $4
WHERE film_id = $1;


