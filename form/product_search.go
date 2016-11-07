package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type ProductSearch struct {
	SearchBase
	ID                    bson.ObjectId
	Name                  string
	QuantityLessThanEqual int
	QuantityMoreThanEqual int
	PriceLessThanEqual    float64
	PriceMoreThanEqual    float64
}

func (p *ProductSearch) ToSearchQuery() (query bson.M) {
	query = bson.M{}
	if p.ID != "" {
		query["_id"] = p.ID
	}

	if p.Name != "" {
		query["name"] = bson.RegEx{
			Pattern: p.Name,
			Options: "i",
		}
	}

	cQuery := bson.M{}
	if p.QuantityLessThanEqual != 0 {
		cQuery["$lte"] = p.QuantityLessThanEqual
	}

	if p.QuantityMoreThanEqual != 0 {
		cQuery["$gte"] = p.QuantityMoreThanEqual
	}

	if len(cQuery) > 0 {
		query["quantity"] = cQuery
	}

	pQuery := bson.M{}
	if p.PriceLessThanEqual != 0 {
		pQuery["$lte"] = p.PriceLessThanEqual
	}

	if p.PriceMoreThanEqual != 0 {
		pQuery["$gte"] = p.PriceMoreThanEqual
	}

	if len(pQuery) > 0 {
		query["price"] = pQuery
	}

	return query
}

func (p *ProductSearch) Validate() (err error) {
	p.Page, p.PerPage, p.OrderBy, err = validateSearchBase(p.Page, p.PerPage, p.OrderBy)

	if p.ID != "" {
		if !bson.IsObjectIdHex(p.ID.Hex()) {
			err = errors.New("Invalid Product.id: not a valid ObjectId")
		}
	}

	if p.Name != "" {
		if !govalidator.IsByteLength(p.Name, 1, 128) {
			err = errors.New("Invalid Product.name: must be between 1 and 128 characters long")
		}
	}

	if p.QuantityLessThanEqual < 0 {
		p.QuantityLessThanEqual = 0
	}
	if p.QuantityMoreThanEqual < 0 {
		p.QuantityMoreThanEqual = 0
	}
	if p.QuantityLessThanEqual > 0 && p.QuantityMoreThanEqual > 0 {
		if p.QuantityLessThanEqual > p.QuantityMoreThanEqual {
			err = errors.New("Invalid Product.quantity_less_than_equal: cannot be more than Product.quantity_more_than_equal")
		}
	}

	if p.PriceLessThanEqual < 0 {
		p.PriceLessThanEqual = 0
	}
	if p.PriceMoreThanEqual < 0 {
		p.PriceMoreThanEqual = 0
	}
	if p.PriceLessThanEqual > 0 && p.PriceMoreThanEqual > 0 {
		if p.PriceLessThanEqual > p.PriceMoreThanEqual {
			err = errors.New("Invalid Product.price_less_than_equal: cannot be more than Product.price_more_than_equal")
		}
	}

	return err
}
