package api

import (
	"6_gopro/product"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

func AddProduct(c echo.Context) error {

	key := []byte("secret")
	claims := jwt.MapClaims{}
	token, err := c.Cookie("token")
	Product := new(product.BindProduct)

	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to read cookie"})
	}

	_, err = jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to parse token"})
	}

	if !claims["isadmin"].(bool) {
		return c.JSON(400, echo.Map{"Error": "You are not access"})
	}

	err = c.Bind(Product)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to bind"})
	}

	err = c.Validate(Product)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unvalid body data"})
	}

	if product.IsProductRepetitive(Product.Name) {
		return c.JSON(400, echo.Map{"Error": "This product is alrealy inserted"})
	}

	pro, err := product.InsertProduct(Product)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Can't insert product"})
	}

	return c.JSON(200, echo.Map{"Status": "Product is add", "Product": pro.Rest()})

}

func getproducts(c echo.Context) error {

	result, err := product.GetAllProducts()
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Can't get products"})
	}

	return c.JSON(200, result)

}
