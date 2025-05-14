package initializations

import (
	"bmt_product_service/internal/routers"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	// Routers
	productRouter := routers.ProductServiceRouterGroup.Product

	mainGroup := r.Group("/v1")
	{
		productRouter.InitProductRouter(mainGroup)
	}

	return r
}
