package api

import (
	"6_gopro/user"

	"github.com/labstack/echo"
)

func Routing(e *echo.Echo) {

	e.POST("/login", Login)
	e.POST("/signup", SignUp)

	group := e.Group("/", user.Auth)

	user := group.Group("user")
	user.GET("/logout", Logout)

	product := group.Group("product")
	product.POST("/add_product", AddProduct)
	product.GET("/get_products", getproducts)

	order := group.Group("order")
	order.GET("/add_order", Order)
	order.GET("/details", OrderDetails)

}
