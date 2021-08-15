package actions

import (
	"doors/models"
	"fmt"
	"log"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_DoorsResource_List() {
	as.LoadFixture("have some doors")
	res := as.HTML("/doors").Get()

	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "This door opens room #1")
	as.Contains(body, "This door opens room #2")
}

func (as *ActionSuite) Test_DoorsResource_Show() {
	as.LoadFixture("have some doors")
	door := &models.Door{}
	err := as.DB.First(door)
	as.NoError(err)
	as.NotZero(door.ID)
	res := as.HTML("/doors/b37e3244-f915-4b7b-81d6-6c8a1edea102").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Door for room #1 has a nice color")
}

func (as *ActionSuite) Test_DoorsResource_Create() {
	as.LoadFixture("have some companies")
	id, _ := uuid.NewV1()
	companyID, _ := uuid.FromString("937af041-43ba-45d0-87a1-6bb173011996")
	log.Printf("companyID: %v", companyID)
	door := &models.Door{
		ID:          id,
		Room:        "#3",
		Floor:       "#3",
		Building:    "#3",
		Description: nulls.NewString("Door won't budge"),
		CompanyID:   companyID,
	}
	res := as.HTML("/doors").Post(door)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/doors/%s", door.ID), res.Location())

	err := as.DB.Eager().First(door)
	as.NoError(err)
	as.NotZero(door.ID)
	as.NotZero(door.CreatedAt)
	as.Equal("#3", door.Room)
	as.Equal("#3", door.Floor)
	as.Equal("#3", door.Building)
	as.Equal(nulls.NewString("Door won't budge"), door.Description)
	as.NotZero(door.Company)
}

func (as *ActionSuite) Test_DoorsResource_Create_Errors() {
	door := &models.Door{
		ID:        uuid.UUID{},
		CompanyID: uuid.UUID{},
	}
	res := as.HTML("/doors").Post(door)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Room can not be blank.")

	c, err := as.DB.Count(door)
	as.NoError(err)
	as.Equal(0, c)
}

func (as *ActionSuite) Test_DoorsResource_Update() {
	as.LoadFixture("have some doors")
	door := &models.Door{}
	err := as.DB.First(door)
	as.NoError(err)
	as.NotZero(door.ID)
	as.NotZero(door.CreatedAt)

	door.Description = nulls.NewString("Lock has been removed for this door.")

	res := as.HTML("/doors/%s", door.ID).Put(door)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/doors/%s", door.ID), res.Location())

	err = as.DB.Reload(door)
	as.NoError(err)
	as.Equal(nulls.NewString("Lock has been removed for this door."), door.Description)
}

func (as *ActionSuite) Test_DoorsResource_Destroy() {
	as.LoadFixture("have some doors")
	door := &models.Door{}
	err := as.DB.First(door)
	as.NoError(err)
	as.NotZero(door.ID)
	as.NotZero(door.CreatedAt)
	res := as.HTML("/doors/b37e3244-f915-4b7b-81d6-6c8a1edea102").Delete()
	as.Equal(303, res.Code)
	as.Equal("/doors", res.Location())
	count, err := as.DB.Count(&models.Doors{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_DoorsResource_New() {
	res := as.HTML("/doors/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New Door")
}

func (as *ActionSuite) Test_DoorsResource_Edit() {
	as.LoadFixture("have some doors")
	door := &models.Door{}
	err := as.DB.First(door)
	as.NoError(err)
	res := as.HTML("/doors/b37e3244-f915-4b7b-81d6-6c8a1edea102").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Door for room #1 has a nice color")
}
