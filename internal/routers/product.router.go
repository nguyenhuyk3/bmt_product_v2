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
		adminFilmRouterPrivate := filmRouter.Group("/admin")
		{
			adminFilmRouterPrivate.POST("/add",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				productController.AddFilm)
			adminFilmRouterPrivate.PUT("/update",
				getFromHeaderMiddleware.GetEmailFromHeader(),
				productController.UpdateFilm)
			adminFilmRouterPrivate.GET("/get_all_films",
				productController.GetAllFilms)
			adminFilmRouterPrivate.POST("/check_and_cache_film_existence/:film_id",
				productController.CheckAndCacheFilmExistence)
		}

		publicFilmRouter := filmRouter.Group("/public")
		{
			publicFilmRouter.GET("/get_film_by_id/:film_id", productController.GetFilmById)
		}
	}

	fabRouter := router.Group("/fab")
	{
		adminFABRouterPrivate := fabRouter.Group("/admin")
		{
			adminFABRouterPrivate.POST("/add", productController.AddFAB)
			adminFABRouterPrivate.PUT("/update", productController.UpdateFAB)
			adminFABRouterPrivate.POST("/delete", productController.DeleteFAB)
		}
	}
}
