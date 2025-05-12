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

	filmRouterPublic := router.Group("/film")
	{
		filmRouterPrivate := filmRouterPublic.Group("/admin")
		{
			filmRouterPrivate.POST("/add", getFromHeaderMiddleware.GetEmailFromHeader(), productController.AddFilm)
			filmRouterPrivate.PUT("/update", getFromHeaderMiddleware.GetEmailFromHeader(), productController.UpdateFilm)
			filmRouterPrivate.GET("/get_all_films", productController.GetAllFilms)
			filmRouterPrivate.POST("check_and_cache_film_existence", productController.CheckAndCacheFilmExistence)
		}
	}

	fabRouterPublic := router.Group("/fab")
	{
		fabRouterPrivate := fabRouterPublic.Group("/admin")
		{
			fabRouterPrivate.POST("/add", productController.AddFAB)
			fabRouterPrivate.PUT("/update", productController.UpdateFAB)
			fabRouterPrivate.POST("/delete/:id", productController.DeleteFAB)
		}
	}
}
