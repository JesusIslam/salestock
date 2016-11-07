package form

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CouponSearch struct {
	SearchBase
	ID                    bson.ObjectId
	QuantityLessThanEqual int
	QuantityMoreThanEqual int
	ValidUntilBefore      time.Time
	ValidUntilAfter       time.Time
}

func (c *CouponSearch) ToSearchQuery() (query bson.M) {
	if c.ID != "" {
		query["_id"] = c.ID
	}

	qQuery := bson.M{}
	if c.QuantityLessThanEqual != 0 {
		qQuery["$lte"] = c.QuantityLessThanEqual
	}

	if c.QuantityMoreThanEqual != 0 {
		qQuery["$gte"] = c.QuantityMoreThanEqual
	}

	if len(qQuery) > 0 {
		query["quantity"] = qQuery
	}

	vQuery := bson.M{}
	if !c.ValidUntilAfter.IsZero() {
		vQuery["$gte"] = c.ValidUntilAfter
	}

	if !c.ValidUntilBefore.IsZero() {
		vQuery["$lte"] = c.ValidUntilBefore
	}

	if len(vQuery) > 0 {
		query["valid_until"] = vQuery
	}

	return query
}

func (c *CouponSearch) Validate() (err error) {
	c.Page, c.PerPage, c.OrderBy, err = validateSearchBase(c.Page, c.PerPage, c.OrderBy)
	switch c.OrderBy {
	case "-_id":
	case "_id":
	case "quantity":
	case "-quantity":
	case "valid_until":
	case "-valid_until":
		break
	default:
		err = errors.New("Invalid Coupon.order_by: invalid property name")
	}

	if c.ID != "" {
		if !bson.IsObjectIdHex(c.ID.Hex()) {
			err = errors.New("Invalid Coupon.id: not a valid ObjectId")
		}
	}

	if c.QuantityLessThanEqual < 0 {
		c.QuantityLessThanEqual = 0
	}
	if c.QuantityMoreThanEqual < 0 {
		c.QuantityMoreThanEqual = 0
	}
	if c.QuantityLessThanEqual > 0 && c.QuantityMoreThanEqual > 0 {
		if c.QuantityLessThanEqual > c.QuantityMoreThanEqual {
			err = errors.New("Invalid Coupon.quantity_less_than_equal: cannot be more than Coupon.quantity_more_than_equal")
		}
	}

	if !c.ValidUntilBefore.IsZero() && !c.ValidUntilAfter.IsZero() {
		if c.ValidUntilBefore.Before(c.ValidUntilAfter) {
			err = errors.New("Invalid Coupon.valid_until_before: cannot be before Coupon.valid_until_after")
		}
	}

	return err
}
