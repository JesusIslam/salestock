package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type TransactionUpdate struct {
	ID          bson.ObjectId   `json:"id"`
	CouponID    bson.ObjectId   `json:"coupon_id"`
	CustomerID  bson.ObjectId   `json:"customer_id"`
	Products    []bson.ObjectId `json:"products"`
	OrderStatus string          `json:"order_status"`
	Shipment    *Shipment       `json:"shipment"`
}

func (t *TransactionUpdate) ToUpdateData() (id bson.ObjectId, data bson.M) {
	data = bson.M{}
	id = t.ID

	update := bson.M{}
	if t.CouponID != "" {
		update["coupon_id"] = t.CouponID
	}

	if t.CustomerID != "" {
		update["customer_id"] = t.CustomerID
	}

	if t.OrderStatus != "" {
		update["order_status"] = t.OrderStatus
	}

	if t.Shipment != nil {
		shipment := bson.M{}
		if t.Shipment.ID != "" {
			shipment["shipment_id"] = t.Shipment.ID
		}

		if t.Shipment.Status != "" {
			shipment["status"] = t.Shipment.Status
		}

		if t.Shipment.Name != "" {
			shipment["name"] = t.Shipment.Name
		}

		if t.Shipment.PhoneNumber != "" {
			shipment["phone_number"] = t.Shipment.PhoneNumber
		}

		if t.Shipment.Email != "" {
			shipment["email"] = t.Shipment.Email
		}

		if t.Shipment.Address != "" {
			shipment["address"] = t.Shipment.Address
		}

		update["shipment"] = shipment
	}

	data["$set"] = update

	return id, data
}

func (t *TransactionUpdate) Validate() (err error) {
	if !bson.IsObjectIdHex(t.ID.Hex()) {
		err = errors.New("Invalid Transaction.id: not a valid ObjectId")
	}

	if t.CouponID != "" {
		if !bson.IsObjectIdHex(t.CouponID.Hex()) {
			err = errors.New("Invalid Transaction.coupon_id: not a valid ObjectId")
		}
	}

	if t.CustomerID != "" {
		if !bson.IsObjectIdHex(t.CustomerID.Hex()) {
			err = errors.New("Invalid Transaction.customer_id: not a valid ObjectId")
		}
	}

	if len(t.Products) > 0 {
		for _, p := range t.Products {
			if !bson.IsObjectIdHex(p.Hex()) {
				err = errors.New("Invalid Transaction.products: not a valid ObjectId in one of the elements")
			}
		}
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

	if t.Shipment != nil {
		if t.Shipment.ID != "" {
			if !govalidator.IsByteLength(t.Shipment.ID, 1, 128) {
				err = errors.New("Invalid Transaction.shipment.id: must be between 1 and 128 characters long")
			}
		}

		if t.Shipment.Status != "" {
			switch t.Shipment.Status {
			case "not_send":
			case "valid":
			case "invalid":
				break
			default:
				err = errors.New("Invalid Transaction.shipment.status: must be not_send, valid, or invalid")
			}
		}

		if t.Shipment.Name != "" {
			if !govalidator.IsByteLength(t.Shipment.Name, 1, 128) {
				err = errors.New("Invalid Transaction.shipment.name: must be between 1 and 128 characters long")
			}
		}

		if t.Shipment.PhoneNumber != "" {
			if !govalidator.IsByteLength(t.Shipment.PhoneNumber, 6, 24) {
				err = errors.New("Invalid Transaction.shipment.phone_number: must be between 6 and 24 characters long")
			}
		}

		if t.Shipment.Email != "" {
			if !govalidator.IsEmail(t.Shipment.Email) {
				err = errors.New("Invalid Transaction.shipment.email: not a valid email address")
			}
		}

		if t.Shipment.Address != "" {
			if !govalidator.IsByteLength(t.Shipment.Address, 1, 128) {
				err = errors.New("Invalid Transaction.shipment.address: must be between 1 and 128 characters long")
			}
		}
	}

	return err
}
