package order

import (
	"6_gopro/db"
	"6_gopro/product"
	"6_gopro/user"
	"context"
	"errors"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	p = "product"
	u = "user"
	o = "order"
)

type (
	Order struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
		Owner   string             `json:"owner" bson:"owner"`
		Product string             `json:"product" bson:"product"`
		Price   int                `json:"price" bson:"price"`
	}

	Details struct {
		ID      string           `json:"_id" `
		Owner   string           `json:"owner" `
		Product *product.Product `json:"product" `
	}
)

func (d *Details) Rest() echo.Map {
	return echo.Map{
		"Order Id":         d.ID,
		"Order Owner Name": d.Owner,
		"Order Product":    d.Product.Rest(),
	}

}

func (o *Order) Rest() echo.Map {
	return echo.Map{
		"Order Id":         o.ID.Hex(),
		"Order Owner Id":   o.Owner,
		"Order Product Id": o.Product,
		"Order Price":      o.Price,
	}

}

func GetProduct(productId string) (*product.Product, error) {

	result := new(product.Product)

	objID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = db.Db.Collection(p).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func InsertOrder(order *Order) error {

	_, err := db.Db.Collection(o).InsertOne(context.Background(), order)
	if err != nil {
		return errors.New(" Unable to insert order ")
	}
	return nil
}

func GetOrder(id string) (*Order, error) {

	result := new(Order)
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	err := db.Db.Collection(o).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUsername(id string) (string, error) {

	result := new(user.User)
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	err := db.Db.Collection(u).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return "", err
	}
	return result.Username, nil
}
