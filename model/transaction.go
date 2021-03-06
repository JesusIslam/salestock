package model

import (
	"errors"

	"gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	CouponID   bson.ObjectId `bson:"coupon_id,omitempty" json:"coupon_id,omitempty"`
	CustomerID bson.ObjectId `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	Products   []*struct {
		ID        bson.ObjectId `json:"id,omitempty"`
		Name      string        `json:"name,omitempty"`
		PriceEach float64       `json:"price_each,omitempty"`
	} `bson:"products,omitempty" json:"products,omitempty"`
	TotalPrice  float64   `json:"total_price"`
	OrderStatus string    `bson:"order_status,omitempty" json:"order_status,omitempty"` // submitted, valid, invalid
	Shipment    *Shipment `bson:"shipment,omitempty" json:"shipment,omitempty"`
}

func (t *Transaction) Validate() (err error) {
	if t.CouponID != "" {
		if !bson.IsObjectIdHex(t.CouponID.Hex()) {
			err = errors.New("Invalid Transaction.coupon_id: not a valid ObjectId")
		}
	}

	if !bson.IsObjectIdHex(t.CustomerID.Hex()) {
		err = errors.New("Invalid Transaction.customer_id: not a valid ObjectId")
	}

	for _, p := range t.Products {
		if !bson.IsObjectIdHex(p.ID.Hex()) {
			err = errors.New("Invalid Transaction.products element: id not a valid ObjectId found")
		}
		if !govalidator.IsByteLength(p.Name, 1, 128) {
			err = errors.New("Invalid Transaction.products element: name must be between 1 and 128 characters long")
		}
		if p.PriceEach < 0 {
			p.PriceEach = 0
		}
	}

	switch t.OrderStatus {
	case "submitted":
	case "valid":
	case "invalid":
		break
	default:
		err = errors.New("Invalid Transaction.order_status: must be submitted, valid, or invalid")
	}

	if t.Shipment != nil {
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
	ID          string `bson:"id,omitempty" json:"id,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Email       string `bson:"email,omitempty" json:"email,omitempty"`
	Address     string `bson:"address,omitempty" json:"address,omitempty"`
	Status      string `bson:"status,omitempty" json:"status,omitempty"` // not_send, on_progress, arrived, accepted
}
