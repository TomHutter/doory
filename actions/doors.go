package actions

import (
	"doors/models"
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Door)
// DB Table: Plural (doors)
// Resource: Plural (Doors)
// Path: Plural (/doors)
// View Template Folder: Plural (/templates/doors/)

// DoorsResource is the resource for the Door model
type DoorsResource struct {
	buffalo.Resource
}

// List gets all Doors. This function is mapped to the path
// GET /doors
func (v DoorsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	doors := &models.Doors{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Doors from the DB
	if err := q.Eager().All(doors); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("doors", doors)
		return c.Render(http.StatusOK, r.HTML("/doors/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(doors))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(doors))
	}).Respond(c)
}

// Show gets the data for one Door. This function is mapped to
// the path GET /doors/{door_id}
func (v DoorsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Door
	door := &models.Door{}

	// To find the Door the parameter door_id is used.
	if err := tx.Eager().Find(door, c.Param("door_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("door", door)

		return c.Render(http.StatusOK, r.HTML("/doors/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(door))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(door))
	}).Respond(c)
}

// New renders the form for creating a new Door.
// This function is mapped to the path GET /doors/new
func (v DoorsResource) New(c buffalo.Context) error {
	c.Set("door", &models.Door{})

	set_companies(c)
	if err := pushBreadcrumb(c, "New door"); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.HTML("/doors/new.plush.html"))
}

// Create adds a Door to the DB. This function is mapped to the
// path POST /doors
func (v DoorsResource) Create(c buffalo.Context) error {
	// Allocate an empty Door
	door := &models.Door{}

	// Bind door to the html form elements
	if err := c.Bind(door); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	if err := set_companies(c); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	setBreadcrumbs(c)

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(door)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("door", door)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/doors/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "door.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/doors/%v", door.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(door))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(door))
	}).Respond(c)
}

// Edit renders a edit form for a Door. This function is
// mapped to the path GET /doors/{door_id}/edit
func (v DoorsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Door
	door := &models.Door{}

	if err := tx.Eager().Find(door, c.Param("door_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("door", door)

	set_companies(c)

	// Set helper to redirect to previous page
	//set_referrerPath(c)

	label := fmt.Sprintf("Door %s %s", door.Building, door.Room)
	if err := pushBreadcrumb(c, label); err != nil {
		return err
	}

	return c.Render(http.StatusOK, r.HTML("/doors/edit.plush.html"))
}

// Update changes a Door in the DB. This function is mapped to
// the path PUT /doors/{door_id}
func (v DoorsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Door
	door := &models.Door{}

	if err := tx.Find(door, c.Param("door_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Door to the html form elements
	if err := c.Bind(door); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(door)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("door", door)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/doors/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "door.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/doors/%v", door.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(door))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(door))
	}).Respond(c)
}

// Destroy deletes a Door from the DB. This function is mapped
// to the path DELETE /doors/{door_id}
func (v DoorsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Door
	door := &models.Door{}

	// To find the Door the parameter door_id is used.
	if err := tx.Find(door, c.Param("door_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Allocate an empty AccessGroupDoors
	accessGroupDoors := &models.AccessGroupDoors{}

	// Get all AccessGroupDoors belonging to door
	if err := tx.Where("door_id = ?", c.Param("door_id")).All(accessGroupDoors); err != nil {
		return err
	}

	if err := tx.Destroy(accessGroupDoors); err != nil {
		return err
	}

	if err := tx.Destroy(door); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "door.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/doors")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(door))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(door))
	}).Respond(c)
}
