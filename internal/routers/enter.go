package routers

type RouterGroup struct {
	Product ProductRouter
}

var ProductServiceRouterGroup = new(RouterGroup)
