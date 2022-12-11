package api

import (
	"6_gopro/order"
	"6_gopro/user"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Order(c echo.Context) error {

	Order := new(order.Order)

	productId := c.QueryParam("id")

	p, err := order.GetProduct(productId)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Product not found"})
	}

	username, err := user.GetUserFromToken(c)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to get user from token"})
	}

	Order.ID = primitive.NewObjectID()
	Order.Product = p.ID.Hex()
	Order.Price = p.Price
	Order.Owner = username.ID.Hex()

	err = order.InsertOrder(Order)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to instert order"})
	}
	return c.JSON(200, echo.Map{"Status": "Your order has been made", "Order": Order.Rest()})

}

func OrderDetails(c echo.Context) error {

	id := c.QueryParam("id")

	d := new(order.Details)

	Order, err := order.GetOrder(id)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to find order"})
	}

	pro, err := order.GetProduct(Order.Product)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Product not found"})
	}

	username, err := order.GetUsername(Order.Owner)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Owner not found"})
	}

	d.ID = Order.ID.Hex()
	d.Owner = username
	d.Product = pro

	return c.JSON(200, echo.Map{"Order": d.Rest()})

}
