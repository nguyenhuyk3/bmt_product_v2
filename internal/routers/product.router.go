package routers

import (
	"bmt_product_service/internal/injectors"
	"bmt_product_service/internal/middlewares"
	"log"

	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (pr *ProductRouter) InitProductRouter(router *gin.RouterGroup) {
	productController, err := injectors.InitProductController()
	if err != nil {
		log.Fatalf("failed to initialize ProductController: %v", err)
		return
	}

	getFromHeaderMiddleware := middlewares.NewGetFromHeaderMiddleware()

	filmRouter := router.Group("/film")
	{
		privateFilmRouter := filmRouter.Group("/admin")
		{
			privateFilmRouter.POST("/add",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				productController.AddFilm)
			privateFilmRouter.PUT("/update",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				productController.UpdateFilm)
			privateFilmRouter.GET("/get_all_films",
				productController.GetAllFilms)
			privateFilmRouter.POST("/check_and_cache_film_existence/:film_id",
				productController.CheckAndCacheFilmExistence)
		}

		publicFilmRouter := filmRouter.Group("/public")
		{
			publicFilmRouter.GET("/get_film_by_id/:film_id", productController.GetFilmById)
		}
	}

	fabRouter := router.Group("/fab")
	{
		privateFABRouter := fabRouter.Group("/admin")
		{
			privateFABRouter.POST("/add", productController.AddFAB)
			privateFABRouter.PUT("/update", productController.UpdateFAB)
			privateFABRouter.POST("/delete", productController.DeleteFAB)
		}

		publicFABRouter := fabRouter.Group("/public")
		{
			publicFABRouter.GET("/get_all", productController.GetAllFABs)
		}
	}
}
