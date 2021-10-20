package actions

import (
	"doory/models"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_AccessGroupsResource_List() {
	as.LoadFixture("have some access_groups")
	res := as.HTML("/access_groups/").Get()

	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Company #1")
	as.Contains(body, "Company #2")
}

func (as *ActionSuite) Test_AccessGroupsResource_Show() {
	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err := as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	res := as.HTML("/access_groups/30f355b5-7395-4bee-befe-1e6336e1cd4e").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Company #2")
}

func (as *ActionSuite) Test_AccessGroupsResource_Create() {
	id, _ := uuid.NewV4()
	access_group := &models.AccessGroup{
		ID:          id,
		Name:        "Garage",
		Description: nulls.NewString("We live in the shadows"),
	}
	res := as.HTML("/access_groups/").Post(access_group)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/access_groups/%s/edit", access_group.ID), res.Location())

	err := as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	as.NotZero(access_group.Name)
}

func (as *ActionSuite) Test_AccessGroupsResource_Create_Errors() {
	access_group := &models.AccessGroup{
		ID: uuid.UUID{},
	}
	res := as.HTML("/access_groups/").Post(access_group)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Name can not be blank.")

	c, err := as.DB.Count(access_group)
	as.NoError(err)
	as.Equal(0, c)
}

// Create duplicate AccessGroup with same Name leads to redirect_edit

func (as *ActionSuite) Test_AccessGroupsResource_Create_Duplicates() {
	as.LoadFixture("have some access_groups")
	id, _ := uuid.NewV4()
	access_group := &models.AccessGroup{
		ID: id,
		// Name: Common is already taken in fixtures, but unique match is case insensitive
		Name: "common",
	}
	res := as.HTML("/access_groups/").Post(access_group)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Name is already taken")
}

func (as *ActionSuite) Test_AccessGroupsResource_Update() {
	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err := as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	as.NotZero(access_group.CreatedAt)

	access_group.Name = "Heaven"

	res := as.HTML("/access_groups/%s", access_group.ID).Put(access_group)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/access_groups/%s", access_group.ID), res.Location())

	err = as.DB.Reload(access_group)
	as.NoError(err)
	as.Equal("Heaven", access_group.Name)
}

func (as *ActionSuite) Test_AccessGroupsResource_Destroy() {
	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err := as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	as.NotZero(access_group.CreatedAt)
	res := as.HTML("/access_groups/8910574a-92eb-4619-a273-9f50f9b2948d/").Delete()
	as.Equal(303, res.Code)
	as.Equal("/access_groups", res.Location())
	count, err := as.DB.Count(&models.AccessGroups{})
	as.NoError(err)
	as.Equal(2, count)
}

func (as *ActionSuite) Test_AccessGroupsResource_New() {
	res := as.HTML("/access_groups/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New AccessGroup")
}

func (as *ActionSuite) Test_AccessGroupsResource_Edit() {
	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err := as.DB.First(access_group)
	as.NoError(err)
	res := as.HTML("/access_groups/30f355b5-7395-4bee-befe-1e6336e1cd4e/edit").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Company #2")
}
