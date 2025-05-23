-- name: InsertFAB :one
INSERT INTO foods_and_beverages (name, type, image_url, price)
VALUES ($1, $2, $3, $4)
RETURNING id;