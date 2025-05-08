package provider

import (
	"bmt_product_service/db/sqlc"
	"bmt_product_service/global"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProvidePgxPool() *pgxpool.Pool {
	return global.Postgresql
}

func ProvideQueries() *sqlc.Queries {
	return sqlc.New(global.Postgresql)
}
