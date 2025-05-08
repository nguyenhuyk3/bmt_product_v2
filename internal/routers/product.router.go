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

	productRouterPublic := router.Group("/film")
	{
		filmRouterPrivate := productRouterPublic.Group("/admin")
		{
			filmRouterPrivate.POST("/add", getFromHeaderMiddleware.GetEmailFromHeader(), productController.AddFilm)
			filmRouterPrivate.GET("/get_all_films", productController.GetAllFilms)
		}
	}
}
