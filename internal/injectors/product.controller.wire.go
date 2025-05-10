//go:build wireinject

package injectors

import (
	"bmt_product_service/internal/controllers"
	"bmt_product_service/internal/implementations/product"

	"github.com/google/wire"
)

func InitProductController() (*controllers.ProductController, error) {
	wire.Build(
		uploadServiceSet,
		dbSet,
		redisSet,

		product.NewFilmService,
		controllers.NewProductController,
	)

	return &controllers.ProductController{}, nil
}
