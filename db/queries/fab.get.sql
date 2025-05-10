-- name: GetFABById :one
SELECT * FROM food_and_beverage
WHERE id = $1;

-- name: ListFAB :many
SELECT * FROM food_and_beverage
ORDER BY created_at DESC;

-- name: GetFABImageURLByID :one
SELECT image_url
FROM food_and_beverage
WHERE id = $1;

