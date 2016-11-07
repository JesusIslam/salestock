package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type ProductUpdate struct {
	ID       bson.ObjectId `json:"id"`
	Name     string        `json:"string"`
	Price    float64       `json:"price"`
	Quantity int           `json:"quantity"`
}

func (p *ProductUpdate) ToUpdateData() (id bson.ObjectId, data bson.M) {
	data = bson.M{}
	id = p.ID

	update := bson.M{}
	if p.Name != "" {
		update["name"] = p.Name
	}

	if p.Price > 0 {
		update["price"] = p.Price
	}

	data["$set"] = update
	data["$inc"] = bson.M{
		"quantity": p.Quantity,
	}

	return id, data
}

func (p *ProductUpdate) Validate() (err error) {
	if !bson.IsObjectIdHex(p.ID.Hex()) {
		err = errors.New("Invalid Coupon.id: not a valid ObjectId")
	}

	if p.Name != "" {
		if !govalidator.IsByteLength(p.Name, 1, 128) {
			err = errors.New("Invalid Product.name: must be between 1 and 128 characters long")
		}
	}

	if p.Price < 0.0 {
		err = errors.New("Invalid Product.price: cannot be lower than 0.0")
	}

	return err
}
