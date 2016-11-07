package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

type UserSearch struct {
	SearchBase
	ID       bson.ObjectId
	Username string
	Role     string
}

func (u *UserSearch) ToSearchQuery() (query bson.M) {
	query["_id"] = u.ID

	if u.Username != "" {
		query["username"] = bson.RegEx{
			Pattern: u.Username,
			Options: "i",
		}
	}

	if u.Role != "" {
		query["role"] = u.Role
	}

	return query
}

func (u *UserSearch) Validate() (err error) {
	u.Page, u.PerPage, u.OrderBy, err = validateSearchBase(u.Page, u.PerPage, u.OrderBy)

	switch u.OrderBy {
	case "username":
	case "-username":
	case "_id":
	case "-_id":
	case "role":
	case "-role":
		break
	default:
		err = errors.New("Invalid User.order_by: invalid property name")
	}

	if u.ID != "" {
		if !bson.IsObjectIdHex(u.ID.Hex()) {
			err = errors.New("Invalid User.id: not a valid ObjectId")
		}
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
