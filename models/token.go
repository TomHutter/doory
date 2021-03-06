package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Token is used by pop to map your tokens database table to your go code.
type Token struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	TokenID      string        `json:"token_id" db:"token_id"`
	PersonID     uuid.UUID     `json:"person_id" db:"person_id"`
	Person       Person        `json:"person,omitempty" belongs_to:"person"`
	AccessGroups []AccessGroup `json:"access_groups" many_to_many:"token_access_groups" db:"-"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t Token) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tokens is not required by pop and may be deleted
type Tokens []Token

// String is not required by pop and may be deleted
func (t Tokens) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Token) Validate(tx *pop.Connection) (*validate.Errors, error) {
	token := &Token{}

	var count int
	var err error

	count, err = tx.Where("UPPER(token_id) = UPPER(?)", t.TokenID).Count(token)
	if err != nil {
		errors := validate.NewErrors()
		errors.Add("token_id", "error during db lookup tokens-TokenID")
		return errors, err
	}

	if count > 0 {
		if err := tx.Eager().Where("UPPER(token_id) = UPPER(?)", t.TokenID).First(token); err != nil {
			return nil, err
		}
		errors := validate.NewErrors()
		errors.Add("token_id", fmt.Sprintf("TokenID \"%s\" is already in use by %s %s", t.TokenID, token.Person.Name, token.Person.Surname))
		return errors, nil
	}

	return validate.Validate(
		&validators.StringIsPresent{Field: t.TokenID, Name: "TokenID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Token) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Token) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (t *Token) Reverse() string {
	var reversed string
	chars := []rune(t.TokenID)
	for i := len(chars); i != 0; i-- {
		//TODO
	}
	return reversed
}
