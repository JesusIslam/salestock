package model

import (
	"errors"

	"gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Username string        `bson:"username,omitempty" json:"username,omitempty"`
	Status   string        `bson:"status,omitempty" json:"status,omitempty"`
}

func (u *User) Validate() (err error) {
	if !govalidator.IsByteLength(u.Username, 6, 32) {
		err = errors.New("Invalid User.username: must be between 6 and 32 characters long")
	}

	switch u.Status {
	case "admin":
	case "customer":
		break
	default:
		err = errors.New("Invalid User.customer: must be either admin or customer")
	}

	return err
}
