package actions

import (
	"doors/models"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"html/template"
	"strings"
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

func set_company_people(c buffalo.Context, company *models.Company) error {
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

func set_person_tokens(c buffalo.Context, person *models.Person) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	tokens := &models.Tokens{}

	// Retrieve all Tokens from the DB
	if err := tx.Eager().Order("token_id").Where("person_id in (?)", c.Param("person_id")).All(tokens); err != nil {
		return err
	}

	c.Set("tokens", tokens)
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

func set_person_helpers(c buffalo.Context) {
	// Set helper for checkbox active
	c.Set("activeHelper", func(isActive bool) string {
		if isActive {
			return "checked=\"\""
		} else {
			return ""
		}
	})

	// Set helper for checkbox alarm
	c.Set("alarmHelper", func(alarm bool) string {
		if alarm {
			return "checked=\"\""
		} else {
			return ""
		}
	})

	// Set helper for person already created
	c.Set("personExists", func(person *models.Person) bool {
		if person.ID == uuid.Nil {
			return false
		} else {
			return true
		}
	})

	// Show linked names of tokenAccessGroups
	c.Set("accessGroupList", func(token models.Token) template.HTML {
		list := make([]string, 0)
		for _, ag := range token.AccessGroups {
			list = append(list, fmt.Sprintf("<a href=\"/access_groups/%s/\">%s</a>", ag.ID.String(), ag.Name))
		}
		return template.HTML(strings.Join(list, ", "))
	})
}

func set_doors(c buffalo.Context, doors *models.Doors) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Doors from the DB
	if err := q.Order("building,floor,room").All(doors); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	c.Set("doors", doors)
	return nil
}

func accessGroupOpensDoor(ag *models.AccessGroup, d *models.Door) bool {
	for _, agd := range ag.Doors {
		if agd.ID == d.ID {
			return true
		}
	}
	return false
}

func set_opening_doors(c buffalo.Context, doors *models.Doors, openingDoors map[uuid.UUID]bool) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty AccessGroup
	accessGroup := &models.AccessGroup{}

	if c.Param("access_group_id") != "" {
		// Retrieve AccessGroup from the DB
		if err := tx.Eager().Find(accessGroup, c.Param("access_group_id")); err != nil {
			return err
		}

		for _, door := range *doors {
			name := fmt.Sprintf("Door-%s", door.ID)
			if c.Param(name) == "true" || accessGroupOpensDoor(accessGroup, &door) {
				openingDoors[door.ID] = true
			} else {
				openingDoors[door.ID] = false
			}
		}
	}

	c.Set("openingDoors", openingDoors)

	// Ceck box if accessGroup can open door
	c.Set("opensDoorHelper", func(door models.Door) string {
		if openingDoors[door.ID] {
			return "checked=\"\""
		} else {
			return ""
		}
	})

	return nil
}

func set_access_groups(c buffalo.Context, accessGroups *models.AccessGroups) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all AccessGroups from the DB
	if err := q.Order("name").All(accessGroups); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	c.Set("accessGroups", accessGroups)
	return nil
}

func set_used_access_groups(c buffalo.Context, accessGroups *models.AccessGroups, usedAccessGroups map[uuid.UUID]bool) error {

	for _, ag := range *accessGroups {
		name := fmt.Sprintf("AccessGroup-%s", ag.ID)
		if c.Param(name) == "true" {
			usedAccessGroups[ag.ID] = true
		} else {
			usedAccessGroups[ag.ID] = false
		}
	}

	c.Set("usedAccessGroups", usedAccessGroups)

	// Set helper for used access groups
	c.Set("usedAccessGroupHelper", func(ag models.AccessGroup) string {
		if usedAccessGroups[ag.ID] {
			return "checked=\"\""
		} else {
			return ""
		}
	})

	return nil
}

// Redirect to previous page
func redirect_not_implemented(c buffalo.Context, method string) error {
	message := fmt.Sprintf("%s not implemented", method)
	c.Flash().Add("warning", message)

	if c.Request().Referer() != "" {
		return c.Redirect(302, "%s", c.Request().Referer())
	}
	if c.Param("person_id") != "" {
		if c.Param("token_id") != "" {
			return c.Redirect(302, "/people/%s/tokens/%s", c.Param("person_id"), c.Param("token_id"))
		}
		return c.Redirect(302, "/people/%s", c.Param("person_id"))
	}
	return c.Redirect(302, "/people/")
}
