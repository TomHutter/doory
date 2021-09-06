package actions

import (
	"doors/models"
	"fmt"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_AccessGroupDoorsResource_List() {
	res := as.HTML("/access_group_doors/").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_AccessGroupDoorsResource_Show() {
	res := as.HTML("/access_group_doors/30f355b5-7395-4bee-befe-1e6336e1cd4e").Get()
	as.Equal(404, res.Code)
}

// Method Create is used for creating an destroying AccessGroupDoors
// It's controlled by flag Door-<uuid>=true|false
func (as *ActionSuite) Test_AccessGroupDoorsResource_Create() {
	as.LoadFixture("have some doors")
	door := &models.Door{}
	err := as.DB.First(door)
	as.NoError(err)
	as.NotZero(door.ID)
	as.NotZero(door.CreatedAt)

	as.LoadFixture("have some access_groups")
	access_group := &models.AccessGroup{}
	err = as.DB.First(access_group)
	as.NoError(err)
	as.NotZero(access_group.ID)
	as.NotZero(access_group.CreatedAt)

	id, _ := uuid.NewV4()

	// first we create an AccessGroupDoor

	// create map[string]interface{} which behaves like models.AccessGroupDoor
	access_group_door := map[string]interface{}{
		"ID":            id,
		"DoorID":        door.ID,
		"AccessGroupID": access_group.ID,
	}
	access_group_door["door_id"] = door.ID
	access_group_door["access_group_id"] = access_group.ID
	access_group_door[fmt.Sprintf("Door-%s", door.ID)] = "true"
	res := as.HTML("/access_group_doors/").Post(access_group_door)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/access_groups/%s/edit/", access_group.ID), res.Location())

	agd := &models.AccessGroupDoor{}
	err = as.DB.First(agd)
	as.NoError(err)
	as.NotZero(agd.ID)
	as.NotZero(agd.DoorID)
	as.NotZero(agd.AccessGroupID)

	// and then we destroy it again

	// create map[string]interface{} which behaves like models.AccessGroupDoor
	access_group_door[fmt.Sprintf("Door-%s", door.ID)] = "false"
	res = as.HTML("/access_group_doors/").Post(access_group_door)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/access_groups/%s/edit/", access_group.ID), res.Location())

	count, err := as.DB.Count(&models.AccessGroupDoors{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_AccessGroupDoorsResource_Update() {
	as.LoadFixture("have some access_group_doors")
	access_group_door := &models.AccessGroupDoor{}
	err := as.DB.First(access_group_door)
	as.NoError(err)
	as.NotZero(access_group_door.ID)
	as.NotZero(access_group_door.CreatedAt)

	id, _ := uuid.NewV4()
	access_group_door.DoorID = id

	res := as.HTML("/access_group_doors/%s", access_group_door.ID).Put(access_group_door)
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_AccessGroupDoorsResource_Destroy() {
	as.LoadFixture("have some access_group_doors")
	access_group_door := &models.AccessGroupDoor{}
	err := as.DB.First(access_group_door)
	as.NoError(err)
	as.NotZero(access_group_door.ID)
	as.NotZero(access_group_door.CreatedAt)

	res := as.HTML("/access_group_doors/%s", access_group_door.ID).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_AccessGroupDoorsResource_New() {
	res := as.HTML("/access_group_doors/new").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_AccessGroupDoorsResource_Edit() {
	as.LoadFixture("have some access_group_doors")
	res := as.HTML("/access_group_doors/1440013e-3c84-4790-8619-8f7b07a24760/edit").Get()
	as.Equal(404, res.Code)
}
