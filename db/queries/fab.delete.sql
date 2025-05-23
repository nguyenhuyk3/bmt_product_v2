-- name: ToggleFABDelete :exec
UPDATE foods_and_beverages
SET is_deleted = NOT is_deleted,
    updated_at = NOW()
WHERE id = $1;
