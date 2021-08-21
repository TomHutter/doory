package actions

import (
	"doors/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

func set_companies(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	companies := &models.Companies{}

	// Retrieve all Companies from the DB
	if err := tx.Order("name").All(companies); err != nil {
		return err
	}

	c.Set("companies", companies)
	return nil
}

func set_people(c buffalo.Context, company *models.Company) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	people := &models.People{}

	// Retrieve all Companies from the DB
	if err := tx.Order("surname").Where("company_id in (?)", company.ID).All(people); err != nil {
		return err
	}

	c.Set("people", people)
	return nil
}

func set_person(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Person
	person := &models.Person{}

	// To find the Person the parameter person_id is used.
	if err := tx.Find(person, c.Param("person_id")); err != nil {
		return err
	}

	c.Set("person", person)
	return nil
}
