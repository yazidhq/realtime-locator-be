package model

import (
	"TeamTrackerBE/internal/utils/responses"
	"database/sql/driver"
	"fmt"
)

type (
	Role           string
	ContactStatus  string
)

const (
	Superadmin Role = "superadmin"
	Admin      Role = "admin"
	Customer   Role = "customer"

	Pending  ContactStatus = "pending"
	Accepted ContactStatus = "accepted"
	Declined ContactStatus = "declined"
)

func (r Role) IsValid() bool {
	switch r {
	case Superadmin, Admin, Customer:
		return true
	}
	return false
}

func (r ContactStatus) IsValid() bool {
	switch r {
	case Pending, Accepted, Declined:
		return true
	}
	return false
}

func scanString[T ~string](dest *T, value any) error {
	if value == nil {
		*dest = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*dest = T(v)
	case string:
		*dest = T(v)
	default:
		return responses.NewBadRequestError(fmt.Sprintf("cannot scan value %v (type %T) into %T", value, value, dest))
	}
	return nil
}

func (r *Role) Scan(value any) error {
	if err := scanString(r, value); err != nil {
		return err
	}
	if !r.IsValid() {
		return responses.NewBadRequestError(fmt.Sprintf("invalid Role value: %s", *r))
	}
	return nil
}

func (r *ContactStatus) Scan(value any) error {
	if err := scanString(r, value); err != nil {
		return err
	}
	if !r.IsValid() {
		return responses.NewBadRequestError(fmt.Sprintf("invalid Contact Status value: %s", *r))
	}
	return nil
}

func (r Role) Value() (driver.Value, error) {
	if !r.IsValid() {
		return nil, responses.NewBadRequestError(fmt.Sprintf("invalid Role value: %s", r))
	}
	return string(r), nil
}

func (r ContactStatus) Value() (driver.Value, error) {
	if !r.IsValid() {
		return nil, responses.NewBadRequestError(fmt.Sprintf("invalid Contact Status value: %s", r))
	}
	return string(r), nil
}

var Models = []any{
	&User{},
	&Contact{},
}
