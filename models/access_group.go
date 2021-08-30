package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"time"
)

// AccessGroup is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type AccessGroup struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description nulls.String `json:"description" db:"description"`
	Tokens      []Token      `json:"tokens" many_to_many:"token_access_groups" db:"-"`
	Doors       []Door       `json:"doors" many_to_many:"access_group_doors" db:"-"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a AccessGroup) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AccessGroups is not required by pop and may be deleted
type AccessGroups []AccessGroup

// String is not required by pop and may be deleted
func (a AccessGroups) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AccessGroup) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AccessGroup) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AccessGroup) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// AddDoor adds door to AccessGroup, if door belongs not to Access Group already.
func (a *AccessGroup) AddDoor(tx *pop.Connection, door *Door) (*validate.Errors, error) {
	accessGroupDoors := &AccessGroupDoors{}
	q := tx.RawQuery("SELECT * FROM access_group_doors WHERE access_group_id in (?) and door_id in (?)", a.ID, door.ID)
	var count int
	var err error

	count, err = q.Count(accessGroupDoors)
	if err != nil {
		errors := validate.NewErrors()
		errors.Add("token_id", "error during db lookup tokens-TokenID")
		return errors, err
	}

	// door belongs already to AccessGroup
	if count > 0 {
		return nil, nil
	}

	accessGroupDoor := &AccessGroupDoor{
		AccessGroupID: a.ID,
		DoorID:        door.ID,
	}

	// Validate the data from the html form
	return tx.ValidateAndCreate(accessGroupDoor)
}
