package actions

import (
	"doors/models"
	"fmt"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_PeopleResource_List() {
	as.LoadFixture("have some people")
	res := as.HTML("/people").Get()

	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "luke.skywalker@tatooine.com")
	as.Contains(body, "leia.organa@alderaan.com")
}

func (as *ActionSuite) Test_PeopleResource_Show() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "luke.skywalker@tatooine.com")
}

func (as *ActionSuite) Test_PeopleResource_Create() {
	id, _ := uuid.NewV4()
	cpID, _ := uuid.FromString("33e166d8-ee6f-4017-8e16-cd0742f4e0e2")
	person := &models.Person{
		ID:        id,
		Name:      "Lando",
		Surname:   "Calrissian",
		CompanyID: cpID,
		Phone:     "99887766",
		Email:     "lando.calrissian@soccoro.com",
		IDNumber:  "lca",
	}
	res := as.HTML("/people").Post(person)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/people/%s", person.ID), res.Location())

	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	as.NotZero(person.CreatedAt)
	as.Equal("Lando", person.Name)
	as.Equal("Calrissian", person.Surname)
	as.Equal("99887766", person.Phone)
	as.Equal("lando.calrissian@soccoro.com", person.Email)
	as.NotZero(person.CompanyID)
}

func (as *ActionSuite) Test_PeopleResource_Create_Errors() {
	person := &models.Person{
		ID:        uuid.UUID{},
		CompanyID: uuid.UUID{},
	}
	res := as.HTML("/people").Post(person)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Name can not be blank.")

	c, err := as.DB.Count(person)
	as.NoError(err)
	as.Equal(0, c)
}

func (as *ActionSuite) Test_PeopleResource_Update() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	as.NotZero(person.CreatedAt)

	person.Phone = "99999999"

	res := as.HTML("/people/%s", person.ID).Put(person)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/people/%s", person.ID), res.Location())

	err = as.DB.Reload(person)
	as.NoError(err)
	as.Equal("99999999", person.Phone)
}

// Update PeopleResource and set Redirect to "index"
// Response should redirect to /people
func (as *ActionSuite) Test_PeopleResource_Update_List_View() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	as.NotZero(person.CreatedAt)
	as.True(person.IsActive)
	as.False(person.Alarm)

	// Person is used by pop to map your people database table to your go code.
	type Person struct {
		ID       uuid.UUID
		IsActive bool
		Redirect string
	}

	p := Person{person.ID, false, "index"}

	res := as.HTML("/people/%s", person.ID).Put(p)
	as.Equal(303, res.Code)
	as.Equal("/people/", res.Location())

	err = as.DB.Reload(person)
	as.NoError(err)
	as.Equal(false, person.IsActive)
}

func (as *ActionSuite) Test_PeopleResource_Destroy() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	as.NotZero(person.ID)
	as.NotZero(person.CreatedAt)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d").Delete()
	as.Equal(303, res.Code)
	as.Equal("/people", res.Location())
	count, err := as.DB.Count(&models.People{})
	as.NoError(err)
	as.Equal(3, count)
}

func (as *ActionSuite) Test_PeopleResource_New() {
	res := as.HTML("/people/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New Person")
}

func (as *ActionSuite) Test_PeopleResource_Edit() {
	as.LoadFixture("have some people")
	person := &models.Person{}
	err := as.DB.First(person)
	as.NoError(err)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "luke.skywalker@tatooine.com")
}
