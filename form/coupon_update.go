package form

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CouponUpdate struct {
	ID           bson.ObjectId `json:"id"`
	Quantity     int           `json:"quantity"`
	ValidUntil   time.Time     `json:"valid_until"`
	Discount     float64       `json:"discount"`
	DiscountType string        `json:"discount_type"`
}

func (c *CouponUpdate) ToUpdateData() (id bson.ObjectId, data bson.M) {
	id = c.ID

	data["quantity"] = bson.M{
		"$inc": c.Quantity,
	}

	if !c.ValidUntil.IsZero() {
		data["valid_until"] = c.ValidUntil
	}

	if c.Discount != 0 {
		data["discount"] = c.Discount
	}

	if c.DiscountType != "" {
		data["discount_type"] = c.DiscountType
	}

	return id, data
}

func (c *CouponUpdate) Validate() (err error) {
	if !bson.IsObjectIdHex(c.ID.Hex()) {
		err = errors.New("Invalid Coupon.id: not a valid ObjectId")
	}

	if !c.ValidUntil.IsZero() {
		if time.Now().After(c.ValidUntil) {
			err = errors.New("Invalid Coupon.valid_until: cannot be before current server time")
		}
	}

	if c.Discount < 0 {
		err = errors.New("Invalid Coupon.discount: cannot be less than 0")
	}

	if c.DiscountType != "" {
		switch c.DiscountType {
		case "percentage":
		case "nominal":
			break
		default:
			err = errors.New("Invalid Coupon.discount_type: must be percentage or nominal")
		}
	}

	return err
}
