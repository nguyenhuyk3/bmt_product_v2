-- name: InsertFAB :one
INSERT INTO food_and_beverage (name, type, image_url, price)
VALUES ($1, $2, $3, $4)
RETURNING id;