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

// Company is used by pop to map your companies database table to your go code.
type Company struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	Name            string       `json:"name" db:"name"`
	Description     nulls.String `json:"description" db:"description"`
	ContactPersonID uuid.UUID    `json:"-" db:"contact_person_id"`
	//ContactPerson   ContactPerson    `json:"contact_person" db:"contact_person"`
	Doors     []Door    `json:"doors,omitempty" has_many:"doors" order_by:"name desc"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (c Company) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Companies is not required by pop and may be deleted
type Companies []Company

// String is not required by pop and may be deleted
func (c Companies) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Company) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Company) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Company) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Implement the Selectable interface type which provides two behaviours SelectValue and SelectLabel.
// This will allow any list of Companies fetched from the database to be used in the <select> options.
func (c Company) SelectLabel() string {
	return c.Name
}

func (c Company) SelectValue() interface{} {
	return c.ID
}
