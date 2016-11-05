package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Coupon struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Quantity   int           `bson:"quantity,omitempty" json:"quantity,omitempty"`
	ValidUntil time.Time     `bson:"valid_until,omitempty" json:"valid_until,omitempty"`
}

func (c *Coupon) Validate() (err error) {
	if c.Quantity < 1 {
		c.Quantity = 0
	}

	if time.Now().After(c.ValidUntil) {
		err = errors.New("Invalid Coupon.valid_until: cannot be before current server time")
	}

	return err
}
