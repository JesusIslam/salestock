package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type TransactionCreate struct {
	CouponID    bson.ObjectId   `json:"coupon_id"`
	CustomerID  bson.ObjectId   `json:"customer_id"`
	Products    []bson.ObjectId `json:"products"`
	OrderStatus string          `json:"order_status"`
	Shipment    *Shipment       `json:"shipment"`
}

func (t *TransactionCreate) Validate() (err error) {
	if t.CouponID != "" {
		if !bson.IsObjectIdHex(t.CouponID.Hex()) {
			err = errors.New("Invalid Transaction.coupon_id: not a valid ObjectId")
		}
	}

	if !bson.IsObjectIdHex(t.CustomerID.Hex()) {
		err = errors.New("Invalid Transaction.customer_id: not a valid ObjectId")
	}

	if len(t.Products) > 0 {
		for _, p := range t.Products {
			if !bson.IsObjectIdHex(p.Hex()) {
				err = errors.New("Invalid Transaction.products: not a valid ObjectId in one of the elements")
			}
		}
	} else {
		err = errors.New("Invalid Transaction.products: cannot be empty")
	}

	if t.OrderStatus != "" {
		switch t.OrderStatus {
		case "not_send":
		case "on_progress":
		case "arrived":
		case "accepted":
			break
		default:
			err = errors.New("Invalid Transaction.order_status: must be not_send, on_progress, arrived, or accepted")
		}
	}

	if t.Shipment == nil {
		err = errors.New("Invalid Transaction.shipment: cannot be nil")
	} else {
		if !govalidator.IsByteLength(t.Shipment.Name, 1, 128) {
			err = errors.New("Invalid Transaction.shipment.name: must be between 1 and 128 characters long")
		}

		if !govalidator.IsByteLength(t.Shipment.PhoneNumber, 6, 24) {
			err = errors.New("Invalid Transaction.shipment.phone_number: must be between 6 and 24 characters long")
		}

		if !govalidator.IsEmail(t.Shipment.Email) {
			err = errors.New("Invalid Transaction.shipment.email: not a valid email address")
		}

		if !govalidator.IsByteLength(t.Shipment.Address, 1, 128) {
			err = errors.New("Invalid Transaction.shipment.address: must be between 1 and 128 characters long")
		}
	}

	return err
}

type Shipment struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Address     string `json:"address,omitempty"`
	Status      string `json:"status,omitempty"` // not_send, on_progress, arrived
}
