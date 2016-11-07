package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type TransactionSearch struct {
	SearchBase
	ID          bson.ObjectId
	CouponID    bson.ObjectId
	CustomerID  bson.ObjectId
	ProductID   bson.ObjectId
	OrderStatus string
	ShipmentID  bson.ObjectId
	Name        string
	PhoneNumber string
	Email       string
	Address     string
	Status      string
}

func (t *TransactionSearch) ToSearchQuery() (query bson.M) {
	if t.ID != "" {
		query["_id"] = t.ID
	}

	if t.CouponID != "" {
		query["coupon_id"] = t.CouponID
	}

	if t.CustomerID != "" {
		query["customer_id"] = t.CustomerID
	}

	if t.ProductID != "" {
		query["products"] = bson.M{
			"$in": []bson.ObjectId{t.ProductID},
		}
	}

	if t.OrderStatus != "" {
		query["order_status"] = t.OrderStatus
	}

	if t.ShipmentID != "" {
		query["shipment.id"] = t.ShipmentID
	}

	if t.Status != "" {
		query["shipment.status"] = t.Status
	}

	if t.Name != "" {
		query["shipment.name"] = bson.RegEx{
			Pattern: t.Name,
		}
	}

	if t.PhoneNumber != "" {
		query["shipment.phone_number"] = bson.RegEx{
			Pattern: t.PhoneNumber,
			Options: "i",
		}
	}

	if t.Email != "" {
		query["shipment.email"] = bson.RegEx{
			Pattern: t.Email,
			Options: "i",
		}
	}

	if t.Address != "" {
		query["shipment.address"] = bson.RegEx{
			Pattern: t.Address,
			Options: "i",
		}
	}

	return query
}

func (t *TransactionSearch) Validate() (err error) {
	t.Page, t.PerPage, t.OrderBy, err = validateSearchBase(t.Page, t.PerPage, t.OrderBy)

	switch t.OrderBy {
	case "_id":
	case "-_id":
	case "coupon_id":
	case "-coupon_id":
	case "customer_id":
	case "-customer_id":
	case "product_id":
	case "-product_id":
	case "order_status":
	case "-order_status":
	case "shipment.id":
	case "-shipment.id":
	case "shipment.name":
	case "-shipment.name":
	case "shipment.phone_number":
	case "-shipment.phone_number":
	case "shipment.email":
	case "-shipment.email":
	case "shipment.address":
	case "-shipment.address":
	case "shipment.status":
	case "-shipment.status":
		break
	default:
		err = errors.New("Invalid Transaction.order_by: invalid property name")
	}

	if t.ID != "" {
		if !bson.IsObjectIdHex(t.ID.Hex()) {
			err = errors.New("Invalid Transaction.id: not a valid ObjectId")
		}
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

	if t.ProductID != "" {
		if !bson.IsObjectIdHex(t.ProductID.Hex()) {
			err = errors.New("Invalid Transaction.products: not a valid ObjectId in one of the elements")
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

	if t.ShipmentID != "" {
		if !bson.IsObjectIdHex(t.ShipmentID.Hex()) {
			err = errors.New("Invalid Transaction.shipment_id: not a valid ObjectId")
		}
	}

	if t.Status != "" {
		switch t.Status {
		case "not_send":
		case "valid":
		case "invalid":
			break
		default:
			err = errors.New("Invalid Transaction.status: must be not_send, valid, or invalid")
		}
	}

	if t.Name != "" {
		if !govalidator.IsByteLength(t.Name, 1, 128) {
			err = errors.New("Invalid Transaction.name: must be between 1 and 128 characters long")
		}
	}

	if t.PhoneNumber != "" {
		if !govalidator.IsByteLength(t.PhoneNumber, 6, 24) {
			err = errors.New("Invalid Transaction.phone_number: must be between 6 and 24 characters long")
		}
	}

	if t.Email != "" {
		if !govalidator.IsEmail(t.Email) {
			err = errors.New("Invalid Transaction.email: not a valid email address")
		}
	}

	if t.Address != "" {
		if !govalidator.IsByteLength(t.Address, 1, 128) {
			err = errors.New("Invalid Transaction.address: must be between 1 and 128 characters long")
		}
	}

	return err
}
