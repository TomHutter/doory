package actions

import (
	"doors/models"
	"fmt"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_TokenAccessGroupsResource_List() {
	res := as.HTML("/token_access_groups/").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TokenAccessGroupsResource_Show() {
	res := as.HTML("/token_access_groups/30f355b5-7395-4bee-befe-1e6336e1cd4e").Get()
	as.Equal(404, res.Code)
}

// Method Create is used for creating an destroying TokenAccessGroups
// It's controlled by flag AccessGroup-<uuid>=true|false
func (as *ActionSuite) Test_TokenAccessGroupsResource_Create() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	as.NotZero(person.CreatedAt)

	as.LoadFixture("have some tokens")
	token := &models.Token{}
	err = as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	as.NotZero(token.CreatedAt)

	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err = as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	as.NotZero(access_group.CreatedAt)

	id, _ := uuid.NewV4()

	// first we create an AccessGroupDoor

	// create map[string]interface{} which behaves like models.AccessGroupDoor
	token_access_group := map[string]interface{}{
		"ID":            id,
		"AccessGroupID": access_group.ID,
		"TokenID":       token.ID,
		"Person":        person.ID,
	}
	token_access_group["token_id"] = token.ID
	token_access_group["access_group_id"] = access_group.ID
	token_access_group[fmt.Sprintf("AccessGroup-%s", access_group.ID)] = "true"
	res := as.HTML("/token_access_groups/").Post(token_access_group)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/people/%s/tokens/%s/edit/", person.ID, token.ID), res.Location())

	tag := &models.TokenAccessGroup{}
	err = as.DB.First(tag)
	as.NoError(err)
	as.NotZero(tag.ID)
	as.NotZero(tag.AccessGroupID)
	as.NotZero(tag.TokenID)

	// and then we destroy it again

	// create map[string]interface{} which behaves like models.AccessGroupDoor
	token_access_group[fmt.Sprintf("AccessGroup-%s", access_group.ID)] = "false"
	res = as.HTML("/token_access_groups/").Post(token_access_group)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/people/%s/tokens/%s/edit/", person.ID, token.ID), res.Location())

	count, err := as.DB.Count(&models.TokenAccessGroups{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TokenAccessGroupsResource_Update() {
	as.LoadFixture("have some token_access_groups")
	token_access_group := &models.TokenAccessGroup{}
	err := as.DB.First(token_access_group)
	as.NoError(err)
	as.NotZero(token_access_group.ID)
	as.NotZero(token_access_group.CreatedAt)

	id, _ := uuid.NewV4()
	token_access_group.TokenID = id

	res := as.HTML("/access_group_doors/%s", token_access_group.ID).Put(token_access_group)
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TokenAccessGroupsResource_Destroy() {
	as.LoadFixture("have some token_access_groups")
	token_access_group := &models.TokenAccessGroup{}
	err := as.DB.First(token_access_group)
	as.NoError(err)
	as.NotZero(token_access_group.ID)
	as.NotZero(token_access_group.CreatedAt)

	res := as.HTML("/token_access_groups/%s", token_access_group.ID).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TokenAccessGroupsResource_New() {
	res := as.HTML("/token_access_groups/new").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TokenAccessGroupsResource_Edit() {
	as.LoadFixture("have some token_access_groups")
	res := as.HTML("/token_access_groups/0d2dacdb-e23e-4e9f-b892-f8071ffab038/edit").Get()
	as.Equal(404, res.Code)
}
