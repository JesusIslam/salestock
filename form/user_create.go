package form

import (
	"errors"

	govalidator "gopkg.in/asaskevich/govalidator.v4"
)

type UserCreate struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (u *UserCreate) Validate() (err error) {
	if !govalidator.IsByteLength(u.Username, 6, 32) {
		err = errors.New("Invalid User.username: must be between 6 and 32 characters long")
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
