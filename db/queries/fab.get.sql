-- name: GetFABById :one
SELECT * FROM foods_and_beverages
WHERE id = $1;

-- name: ListFAB :many
SELECT * FROM foods_and_beverages
ORDER BY created_at DESC;

-- name: GetFABImageURLByID :one
SELECT image_url
FROM foods_and_beverages
WHERE id = $1;

-- name: IsFABExist :one
SELECT EXISTS (
    SELECT 1 FROM foods_and_beverages WHERE id = $1
) AS EXISTS;

