package product

import (
	"6_gopro/db"
	"context"
	"errors"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const p = "product"

type (
	Product struct {
		ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
		Name     string             `json:"name" bson:"name" validate:"required"`
		Price    int                `json:"price" bson:"price" validate:"required"`
		MadeIn   string             `json:"madein" bson:"madein" `
		CreateAt time.Time          `json:"create_At" bson:"create_At" `
	}

	BindProduct struct {
		Name   string `json:"name" form:"name" validate:"required"`
		Price  int    `json:"price" form:"price" validate:"required"`
		MadeIn string `json:"madein" form:"madein" `
	}
)

func (p *Product) Rest() echo.Map {
	return echo.Map{
		"Product Id":        p.ID.Hex(),
		"Product Name":      p.Name,
		"Product Price":     p.Price,
		"Product Made in":   p.MadeIn,
		"Product Create at": p.CreateAt.Local(),
	}
}

func IsProductRepetitive(name string) bool {

	result := new(Product)
	filter := bson.M{"name": name}

	err := db.Db.Collection(p).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return false
	}
	return true
}

func InsertProduct(pro *BindProduct) (*Product, error) {

	product := new(Product)

	product.ID = primitive.NewObjectID()
	product.Name = pro.Name
	product.Price = pro.Price
	product.MadeIn = pro.MadeIn
	product.CreateAt = time.Now().Local()

	_, err := db.Db.Collection(p).InsertOne(context.Background(), product)
	if err != nil {
		return nil, errors.New(" Unable to insert product ")
	}
	return product, nil
}

func GetAllProducts() ([]bson.M, error) {

	var result []bson.M
	var products []Product

	cursor, err := db.Db.Collection(p).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		result = append(result, primitive.M(product.Rest()))
	}

	return result, nil
}
