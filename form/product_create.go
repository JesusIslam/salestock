package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
)

type ProductCreate struct {
	Name     string  `json:"string"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (p *ProductCreate) Validate() (err error) {
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
