package actions

import (
	"doors/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
	"net/http"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Token)
// DB Table: Plural (tokens)
// Resource: Plural (Tokens)
// Path: Plural (/tokens)
// View Template Folder: Plural (/templates/tokens/)

// TokensResource is the resource for the Token model
type TokensResource struct {
	buffalo.Resource
}

// List gets all Tokens. This function is mapped to the path
// GET /tokens
func (v TokensResource) List(c buffalo.Context) error {
	c.Flash().Add("warning", "List not implemented")

	if c.Param("person_id") != "" {
		return c.Redirect(302, "/people/%s", c.Param("person_id"))
	}
	return c.Redirect(302, "/people/")
}

// Show gets the data for one Token. This function is mapped to
// the path GET /tokens/{token_id}
func (v TokensResource) Show(c buffalo.Context) error {
	c.Flash().Add("warning", "Show not implemented")

	if c.Param("person_id") != "" {
		return c.Redirect(302, "/people/%s", c.Param("person_id"))
	}
	return c.Redirect(302, "/people/")
}

// New renders the form for creating a new Token.
// This function is mapped to the path GET /tokens/new
func (v TokensResource) New(c buffalo.Context) error {
	// To find the Person the parameter person_id is used.
	if err := set_person(c); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("token", &models.Token{})

	return c.Render(http.StatusOK, r.HTML("/tokens/new.plush.html"))
}

// Create adds a Token to the DB. This function is mapped to the
// path POST /tokens
func (v TokensResource) Create(c buffalo.Context) error {
	// Allocate an empty Token
	token := &models.Token{}

	// Bind token to the html form elements
	if err := c.Bind(token); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(token)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("token", token)

			// To find the Person the parameter person_id is used.
			if err := set_person(c); err != nil {
				return c.Error(http.StatusNotFound, err)
			}

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/tokens/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "token.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/people/%v", c.Param("person_id"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(token))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(token))
	}).Respond(c)
}

// Edit renders a edit form for a Token. This function is
// mapped to the path GET /tokens/{token_id}/edit
func (v TokensResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// To find the Person the parameter person_id is used.
	if err := set_person(c); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Allocate an empty Token
	token := &models.Token{}

	if err := tx.Eager().Find(token, c.Param("token_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	accessGroups := &models.AccessGroups{}

	if err := set_access_groups(c, accessGroups); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Check box if tokenAccessGroup token, accessGroup exists
	c.Set("tokenAccessGroupHelper", func(token *models.Token, accessGroup models.AccessGroup) string {
		tokenAccessGroup := &models.TokenAccessGroup{}
		if err := tx.Where("token_id = ? and access_group_id = ?", token.ID, accessGroup.ID).First(tokenAccessGroup); err != nil {
			return ""
		} else {
			return "checked=\"\""
		}
	})

	// Set helper for form IDs
	c.Set("formID", func(id uuid.UUID) string {
		return fmt.Sprintf("access-group-%s", id.String())
	})

	usedAccessGroups := make(map[uuid.UUID]bool)

	if err := set_used_access_groups(c, accessGroups, usedAccessGroups); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	//if err := set_access_groups(c, token); err != nil {
	//	return c.Error(http.StatusNotFound, err)
	//}

	c.Set("token", token)
	return c.Render(http.StatusOK, r.HTML("/tokens/edit.plush.html"))
}

// Update changes a Token in the DB. This function is mapped to
// the path PUT /tokens/{token_id}
func (v TokensResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Token
	token := &models.Token{}

	if err := tx.Find(token, c.Param("token_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Token to the html form elements
	if err := c.Bind(token); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(token)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("token", token)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/tokens/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "token.updated.success"))

		// and redirect to the show page
		if c.Param("Redirect") == "index" {
			return c.Redirect(http.StatusSeeOther, "/people/")
		} else {
			return c.Redirect(http.StatusSeeOther, "/people/%v", c.Param("person_id"))
		}

	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(token))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(token))
	}).Respond(c)
}

// Destroy deletes a Token from the DB. This function is mapped
// to the path DELETE /tokens/{token_id}
func (v TokensResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Token
	token := &models.Token{}

	// To find the Token the parameter token_id is used.
	if err := tx.Find(token, c.Param("token_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(token); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "token.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/people/%v", c.Param("person_id"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(token))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(token))
	}).Respond(c)
}