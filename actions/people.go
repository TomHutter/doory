package actions

import (
	"doors/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"net/http"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Person)
// DB Table: Plural (people)
// Resource: Plural (People)
// Path: Plural (/people)
// View Template Folder: Plural (/templates/people/)

// PeopleResource is the resource for the Person model
type PeopleResource struct {
	buffalo.Resource
}

// List gets all People. This function is mapped to the path
// GET /people
func (v PeopleResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	people := &models.People{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all People from the DB
	if err := q.Eager().All(people); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("people", people)
		return c.Render(http.StatusOK, r.HTML("/people/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(people))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(people))
	}).Respond(c)
}

// Show gets the data for one Person. This function is mapped to
// the path GET /people/{person_id}
func (v PeopleResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Person
	person := &models.Person{}

	// To find the Person the parameter person_id is used.
	if err := tx.Eager().Find(person, c.Param("person_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("person", person)

		return c.Render(http.StatusOK, r.HTML("/people/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(person))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(person))
	}).Respond(c)
}

// New renders the form for creating a new Person.
// This function is mapped to the path GET /people/new
func (v PeopleResource) New(c buffalo.Context) error {
	c.Set("person", &models.Person{})

	set_companies(c)

	return c.Render(http.StatusOK, r.HTML("/people/new.plush.html"))
}

// Create adds a Person to the DB. This function is mapped to the
// path POST /people
func (v PeopleResource) Create(c buffalo.Context) error {
	// Allocate an empty Person
	person := &models.Person{}

	// Bind person to the html form elements
	if err := c.Bind(person); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	set_companies(c)

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(person)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("person", person)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/people/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "person.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/people/%v", person.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(person))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(person))
	}).Respond(c)
}

// Edit renders a edit form for a Person. This function is
// mapped to the path GET /people/{person_id}/edit
func (v PeopleResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Person
	person := &models.Person{}

	if err := tx.Find(person, c.Param("person_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("person", person)

	set_companies(c)

	return c.Render(http.StatusOK, r.HTML("/people/edit.plush.html"))
}

// Update changes a Person in the DB. This function is mapped to
// the path PUT /people/{person_id}
func (v PeopleResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Person
	person := &models.Person{}

	if err := tx.Find(person, c.Param("person_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Person to the html form elements
	if err := c.Bind(person); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(person)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("person", person)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/people/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "person.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/people/%v", person.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(person))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(person))
	}).Respond(c)
}

// Destroy deletes a Person from the DB. This function is mapped
// to the path DELETE /people/{person_id}
func (v PeopleResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Person
	person := &models.Person{}

	// To find the Person the parameter person_id is used.
	if err := tx.Find(person, c.Param("person_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(person); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "person.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/people")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(person))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(person))
	}).Respond(c)
}
