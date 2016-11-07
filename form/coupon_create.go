package form

import (
	"errors"
	"time"
)

type CouponCreate struct {
	Quantity   int       `json:"quantity"`
	ValidUntil time.Time `json:"valid_until"`
}

func (c *CouponCreate) Validate() (err error) {
	if c.Quantity < 1 {
		c.Quantity = 0
	}

	if time.Now().After(c.ValidUntil) {
		err = errors.New("Invalid Coupon.valid_until: cannot be before current server time")
	}

	return err
}
