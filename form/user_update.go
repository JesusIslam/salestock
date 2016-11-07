package form

import (
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type UserUpdate struct {
	ID       bson.ObjectId `json:"id"`
	Username string        `json:"username"`
	Role     string        `json:"role"`
	Password string        `json:"password"`
}

func (u *UserUpdate) ToUpdateData() (id bson.ObjectId, data bson.M) {
	id = u.ID

	if u.Username != "" {
		data["username"] = u.Username
	}

	if u.Role != "" {
		data["role"] = u.Role
	}

	if u.Password != "" {
		data["password"] = u.Password
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

	if u.Password != "" {
		if !govalidator.IsByteLength(u.Password, 6, 50) {
			err = errors.New("Invalid User.password: must be between 6 and 50 characters long")
		} else {
			raw, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			u.Password = base64.StdEncoding.EncodeToString(raw)
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
