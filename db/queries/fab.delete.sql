-- name: DeleteFAB :exec
UPDATE food_and_beverage
SET is_deleted = true,
    updated_at = NOW()
WHERE id = $1;
