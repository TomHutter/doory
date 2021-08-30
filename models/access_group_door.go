package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"time"
)

// AccessGroupDoor is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type AccessGroupDoor struct {
	ID            uuid.UUID `json:"id" db:"id"`
	AccessGroupID uuid.UUID `json:"access_group_id" db:"access_group_id"`
	DoorID        uuid.UUID `json:"door_id" db:"door_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a AccessGroupDoor) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AccessGroupDoors is not required by pop and may be deleted
type AccessGroupDoors []AccessGroupDoor

// String is not required by pop and may be deleted
func (a AccessGroupDoors) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AccessGroupDoor) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: a.AccessGroupID, Name: "CompanyID"},
		&validators.UUIDIsPresent{Field: a.DoorID, Name: "DoorID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AccessGroupDoor) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AccessGroupDoor) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
