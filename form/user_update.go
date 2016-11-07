package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type UserUpdate struct {
	ID       bson.ObjectId `json:"id"`
	Username string        `json:"username"`
	Role     string        `json:"role"`
}

func (u *UserUpdate) ToUpdateData() (id bson.ObjectId, data bson.M) {
	id = u.ID

	if u.Username != "" {
		data["username"] = u.Username
	}

	if u.Role != "" {
		data["role"] = u.Role
	}

	return id, data
}

func (u *UserUpdate) Validate() (err error) {
	if !bson.IsObjectIdHex(u.ID.Hex()) {
		err = errors.New("Invalid User.id: not a valid ObjectId")
	}

	if u.Username != "" {
		if !govalidator.IsByteLength(u.Username, 6, 32) {
			err = errors.New("Invalid User.username: must be between 6 and 32 characters long")
		}
	}

	if u.Role != "" {
		switch u.Role {
		case "admin":
		case "customer":
			break
		default:
			err = errors.New("Invalid User.role: must be either admin or customer")
		}
	}

	return err
}
