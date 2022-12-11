package api

import (
	"6_gopro/user"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

func Login(c echo.Context) error {

	User := new(user.Login)

	err := c.Bind(User)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to bind"})
	}

	U, err := user.GetUserByName(User)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Can't find user by name"})
	}

	token, err := user.CreateToken(U.Username, U.IsAdmin)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to generate token"})
	}

	cookie := &http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(1 * time.Hour),
	}

	c.SetCookie(cookie)

	return c.JSON(200, echo.Map{"Status": "Login", "User": U.Rest(), "Token": token})

}

func SignUp(c echo.Context) error {

	User := new(user.Signup)

	_, err := c.Cookie("token")
	if err != nil {

		err = c.Bind(User)
		if err != nil {
			return c.JSON(400, echo.Map{"Error": "Unable to bind"})
		}

		err = c.Validate(User)
		if err != nil {
			return c.JSON(400, echo.Map{"Error": "Unvalid body data"})
		}

		ok := user.IsUserRepetitive(User.Username)
		if ok {
			return c.JSON(400, echo.Map{"Error": "This username is already used"})
		}

		User.IsAdmin = false

		u, err := user.InsertUser(User)
		if err != nil {
			return c.JSON(400, echo.Map{"Error": "Can't insert user"})
		}
		return c.JSON(200, echo.Map{"Status": "User has been made", "user": u.Rest()})

	}

	C_User, err := user.GetUserFromToken(c)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Can't get user from token"})
	}
	if !C_User.IsAdmin {
		return c.JSON(400, echo.Map{"Error": "You are not access . you can't signup when you are login"})
	}

	err = c.Bind(User)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to bind"})
	}

	err = c.Validate(User)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unvalid body data"})
	}

	ok := user.IsUserRepetitive(User.Username)
	if ok {
		return c.JSON(400, echo.Map{"Error": "This username is already used"})
	}

	u, err := user.InsertUser(User)
	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Can't insert user"})
	}
	return c.JSON(200, echo.Map{"Status": "User has been made", "User": u.Rest()})

}

func Logout(c echo.Context) error {

	key := []byte("secret")
	claims := jwt.MapClaims{}
	token, err := c.Cookie("token")

	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to read cookie"})
	}

	_, err = jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return c.JSON(400, echo.Map{"Error": "Unable to parse token"})
	}

	user := claims["username"].(string)

	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  -1}

	c.SetCookie(cookie)

	u := fmt.Sprintf("%s you are logout", user)

	return c.JSON(200, echo.Map{"Status": u})
}
