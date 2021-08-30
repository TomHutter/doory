package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"time"
)
// TokenAccessGroup is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type TokenAccessGroup struct {
    ID uuid.UUID `json:"id" db:"id"`
    TokenID uuid.UUID `json:"token_id" db:"token_id"`
    AccessGroupID uuid.UUID `json:"access_group_id" db:"access_group_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t TokenAccessGroup) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TokenAccessGroups is not required by pop and may be deleted
type TokenAccessGroups []TokenAccessGroup

// String is not required by pop and may be deleted
func (t TokenAccessGroups) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *TokenAccessGroup) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *TokenAccessGroup) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *TokenAccessGroup) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
