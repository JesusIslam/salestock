package model

import (
	"errors"

	"gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string        `bson:"name,omitempty" json:"name,omitempty"`
	Price    float64       `bson:"price,omitempty" json:"price,omitempty"`
	Quantity int           `bson:"quantity,omitempty" json:"quantity,omitempty"`
}

func (p *Product) Validate() (err error) {
	if !govalidator.IsByteLength(p.Name, 1, 128) {
		err = errors.New("Invalid Product.name: must be between 1 and 128 characters long")
	}

	if p.Price < 0.0 {
		err = errors.New("Invalid Product.price: cannot be lower than 0.0")
	}

	if p.Quantity < 1 {
		p.Quantity = 0
	}

	return err
}
