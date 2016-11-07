package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Coupon struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Quantity     int           `bson:"quantity,omitempty" json:"quantity,omitempty"`
	ValidUntil   time.Time     `bson:"valid_until,omitempty" json:"valid_until,omitempty"`
	Discount     float64       `bson:"discount,omitempty" json:"discount,omitempty"`           // in percentage or nominal
	DiscountType string        `bson:"discount_type,omitempty" json:"discount_type,omitempty"` // percentage or nominal
}

func (c *Coupon) Validate() (err error) {
	if c.Quantity < 1 {
		c.Quantity = 0
	}

	if time.Now().After(c.ValidUntil) {
		err = errors.New("Invalid Coupon.valid_until: cannot be before current server time")
	}

	if c.Discount < 0 {
		c.Discount = 0
	}

	switch c.DiscountType {
	case "percentage":
	case "nominal":
		break
	default:
		err = errors.New("Invalid Coupon.discount_type: must be percentage or nominal")
	}

	return err
}
