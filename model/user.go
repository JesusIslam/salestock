package model

import (
	"errors"

	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Username string        `bson:"username,omitempty" json:"username,omitempty"`
	Role     string        `bson:"role,omitempty" json:"role,omitempty"`
	Password string        `bson:"password,omitempty" json:"password,omitempty"`
}

func (u *User) Validate() (err error) {
	if !govalidator.IsByteLength(u.Username, 6, 32) {
		err = errors.New("Invalid User.username: must be between 6 and 32 characters long")
	}

	if !govalidator.IsByteLength(u.Password, 6, 50) {
		err = errors.New("Invalid User.password: must be between 6 and 50 characters long")
	} else {
		raw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = base64.StdEncoding.EncodeToString(raw)
	}

	switch u.Role {
	case "admin":
	case "customer":
		break
	default:
		err = errors.New("Invalid User.role: must be either admin or customer")
	}

	return err
}
