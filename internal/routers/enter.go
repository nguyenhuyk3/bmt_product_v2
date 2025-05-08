package routers

type RouterGroup struct {
	Product ProductRouter
}

var UserServiceRouterGroup = new(RouterGroup)
