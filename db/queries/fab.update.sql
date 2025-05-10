-- name: UpdateFAB :exec
UPDATE food_and_beverage
SET name = $2,
    type = $3,
    price = $4,
    updated_at = NOW()
WHERE id = $1 AND is_deleted = false;

-- name: UpdateFABImageURL :exec
UPDATE food_and_beverage
SET image_url = $1,
    updated_at = NOW()
WHERE id = $2;
