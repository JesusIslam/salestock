package form

import (
	"errors"
	"time"
)

type CouponCreate struct {
	Quantity     int       `json:"quantity"`
	ValidUntil   time.Time `json:"valid_until"`
	Discount     float64   `json:"discount"`
	DiscountType string    `json:"discount_type"`
}

func (c *CouponCreate) Validate() (err error) {
	if c.Quantity < 1 {
		c.Quantity = 0
	}

	if time.Now().After(c.ValidUntil) {
		err = errors.New("Invalid Coupon.valid_until: cannot be before current server time")
	}

	if c.Discount < 0 {
		err = errors.New("Invalid Coupon.discount: cannot be less than 0")
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
