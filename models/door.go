package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
	"github.com/gobuffalo/validate/v3/validators"
)
// Door is used by pop to map your doors database table to your go code.
type Door struct {
    ID uuid.UUID `json:"id" db:"id"`
    Room string `json:"room" db:"room"`
    Floor string `json:"floor" db:"floor"`
    Building string `json:"building" db:"building"`
    Description nulls.String `json:"description" db:"description"`
    Company uuid.UUID `json:"company" db:"company"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (d Door) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Doors is not required by pop and may be deleted
type Doors []Door

// String is not required by pop and may be deleted
func (d Doors) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Door) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Room, Name: "Room"},
		&validators.StringIsPresent{Field: d.Floor, Name: "Floor"},
		&validators.StringIsPresent{Field: d.Building, Name: "Building"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Door) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Door) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
