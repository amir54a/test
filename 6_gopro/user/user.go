package user

import (
	"6_gopro/db"

	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	u = "user"
)

type (
	User struct {
		ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
		Username  string             `json:"username" bson:"username" validate:"min=3,max=15,required"`
		Password  string             `json:"password" bson:"password" validate:"required"`
		Gmail     string             `json:"gmail" bson:"gmail" validate:"required,email"`
		IsAdmin   bool               `json:"isadmin" bson:"isadmin" `
		CreateAt  time.Time          `json:"create_At" bson:"create_At" `
		LastLogin time.Time          `json:"last_login" bson:"last_login" `
	}

	Claim struct {
		Gmail   string `json:"gmail"`
		IsAdmin bool   `json:"isadmin"`
		jwt.StandardClaims
	}

	Login struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	Signup struct {
		Username string `json:"username" form:"username" validate:"min=3,max=15,required"`
		Password string `json:"password" form:"password" validate:"required"`
		Gmail    string `json:"gmail" form:"gmail" validate:"required,email"`
		IsAdmin  bool   `json:"isadmin" form:"isadmin"`
	}
)

func (U *User) Rest() echo.Map {
	return echo.Map{
		"User Id":        U.ID.Hex(),
		"User Username":  U.Username,
		"User Gmail":     U.Gmail,
		"User Is Admin":  U.IsAdmin,
		"User Create At": U.CreateAt,
	}
}

func GetUserByName(user *Login) (*User, error) {

	result := new(User)
	filter := bson.M{"username": user.Username}

	err := db.Db.Collection(u).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return nil, errors.New(" Username not found ")

	}

	if !CheckPasswordHash(user.Password, result.Password) {
		return nil, errors.New(" Password is wrong ")
	}

	LastLogin := time.Now()
	result.LastLogin = LastLogin.Local()

	err = SetLastLoginTime(result)
	if err != nil {
		return nil, errors.New(" Unable to update last login time ")
	}

	return result, nil

}

func SetLastLoginTime(user *User) error {

	filter := bson.M{"username": user.Username}

	_, err := db.Db.Collection(u).UpdateOne(context.Background(), filter, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil

}

func IsUserRepetitive(username string) bool {

	result := new(User)
	filter := bson.M{"username": username}

	err := db.Db.Collection(u).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return false
	}
	return true
}

func InsertUser(s *Signup) (*User, error) {

	user := new(User)
	user.ID = primitive.NewObjectID()

	createAt := time.Now()
	user.CreateAt = createAt.Local()

	hashpassword, err := HashPassword(s.Password)
	if err != nil {
		return nil, errors.New(" Unable to hash password")
	}
	user.Password = hashpassword
	user.Gmail = s.Gmail
	user.Username = s.Username
	user.IsAdmin = s.IsAdmin

	_, err = db.Db.Collection(u).InsertOne(context.Background(), user)
	if err != nil {
		return nil, errors.New(" Unable to insert user ")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(username string, isadmin bool) (string, error) {

	key := []byte("secret")

	Claims := jwt.MapClaims{}
	Claims["username"] = username
	Claims["isadmin"] = isadmin
	Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	token, err := at.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil

}

func InsertAdmin() {

	admin := new(User)

	err := db.Db.Collection(u).FindOne(context.Background(), bson.M{"username": "admin"}).Decode(admin)
	if err != nil {

		admin.ID = primitive.NewObjectID()
		admin.Username = "admin"
		admin.Gmail = "admin@gmail.com"
		admin.IsAdmin = true

		createAt := time.Now()
		admin.CreateAt = createAt.Local()

		hashpassword, _ := HashPassword("123")
		admin.Password = hashpassword

		_, err = db.Db.Collection(u).InsertOne(context.Background(), admin)
		if err != nil {
			panic(err)
		}
	}

}

func GetUserFromToken(c echo.Context) (*User, error) {

	key := []byte("secret")

	claims := jwt.MapClaims{}

	Token, err := c.Cookie("token")
	if err != nil {
		return nil, errors.New(" Unable to find cookie")
	}

	_, err = jwt.ParseWithClaims(Token.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, errors.New(" Unable to pars token")
	}

	result := new(User)
	filter := bson.M{"username": claims["username"].(string)}

	err = db.Db.Collection(u).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return nil, errors.New(" Username not found ")

	}
	return result, nil

}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

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

		n := time.Now().Unix()
		Ntime := float64(n)

		Etime := claims["exp"].(float64)
		if Etime-Ntime < 0 {

			return c.JSON(400, echo.Map{"Error": "Your token is expire"})
		}

		result := new(User)
		Result := claims["username"].(string)
		filter := bson.M{"username": Result}

		err = db.Db.Collection(u).FindOne(context.Background(), filter).Decode(result)
		if err != nil {
			return c.JSON(400, echo.Map{"Error": "User not found"})

		}

		return next(c)
	}
}
