-- name: UpdateFAB :exec
UPDATE food_and_beverage
SET name = $2,
    type = $3,
    image_url = $4,
    price = $5,
    updated_at = NOW()
WHERE id = $1 AND is_deleted = false;