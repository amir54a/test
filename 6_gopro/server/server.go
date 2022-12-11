package server

import (
	"6_gopro/api"
	"6_gopro/config"
	"6_gopro/db"
	"6_gopro/user"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func RunServer() {

	e := echo.New()

	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	api.Routing(e)

	port := fmt.Sprintf(":%s", config.Conf.Listen.Port)

	e.Logger.Fatal(e.Start(port))

}

func init() {

	config.ReadConfig()

	db.Connectdb()

	user.InsertAdmin()

}
