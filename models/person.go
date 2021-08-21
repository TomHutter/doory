package models

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"time"
)

// Person is used by pop to map your people database table to your go code.
type Person struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Surname   string    `json:"surname" db:"surname"`
	CompanyID uuid.UUID `json:"company_id" db:"company_id"`
	Company   Company   `json:"company,omitempty" belongs_to:"company"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	IDNumber  string    `json:"id_number" db:"id_number"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Alarm     bool      `json:"alarm" db:"alarm"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (p Person) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// people is not required by pop and may be deleted
type People []Person

// String is not required by pop and may be deleted
func (p People) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Person) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Surname, Name: "Surname"},
		&validators.StringIsPresent{Field: p.Email, Name: "Email"},
		&validators.StringIsPresent{Field: p.Phone, Name: "Phone"},
		&validators.StringIsPresent{Field: p.IDNumber, Name: "IDNumber"},
		&validators.UUIDIsPresent{Field: p.CompanyID, Name: "CompanyID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Person) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Person) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Implement the Selectable interface type which provides two behaviours SelectValue and SelectLabel.
// This will allow any list of Companies fetched from the database to be used in the <select> options.
func (p Person) SelectLabel() string {
	return fmt.Sprintf("%s %s", p.Name, p.Surname)
}

func (p Person) SelectValue() interface{} {
	return p.ID
}
