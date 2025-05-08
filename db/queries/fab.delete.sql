-- name: ToggleFABDelete :exec
UPDATE food_and_beverage
SET is_deleted = NOT is_deleted,
    updated_at = NOW()
WHERE id = $1;
