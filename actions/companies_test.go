package actions

import (
	"doors/models"
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_CompaniesResource_List() {
	as.LoadFixture("have some companies")
	res := as.HTML("/companies").Get()

	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Company #1")
	as.Contains(body, "Company #2")
}

func (as *ActionSuite) Test_CompaniesResource_Show() {
	as.LoadFixture("have some companies")
	company := &models.Company{}
	err := as.DB.First(company)
	as.NoError(err)
	as.NotZero(company.ID)
	res := as.HTML("/companies/937af041-43ba-45d0-87a1-6bb173011996").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Description for company #1")
}

func (as *ActionSuite) Test_CompaniesResource_Create() {
	//as.LoadFixture("have some companies")
	id, _ := uuid.NewV1()
	cpID, _ := uuid.FromString("4edb21a7-fe0c-11eb-97ee-0242ac1f0002")
	company := &models.Company{
		ID:              id,
		Name:            "Company #3",
		Description:     nulls.NewString("This company is a one man show"),
		ContactPersonID: cpID,
	}
	res := as.HTML("/companies").Post(company)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/companies/%s", company.ID), res.Location())

	err := as.DB.First(company)
	as.NoError(err)
	as.NotZero(company.ID)
	as.NotZero(company.CreatedAt)
	as.Equal("Company #3", company.Name)
	as.Equal(nulls.NewString("This company is a one man show"), company.Description)
	as.NotZero(company.ContactPersonID)
}

func (as *ActionSuite) Test_CompaniesResource_Create_Errors() {
	company := &models.Company{
		ID:              uuid.UUID{},
		ContactPersonID: uuid.UUID{},
	}
	res := as.HTML("/companies").Post(company)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "Name can not be blank.")

	c, err := as.DB.Count(company)
	as.NoError(err)
	as.Equal(0, c)
}

func (as *ActionSuite) Test_CompaniesResource_Update() {
	as.LoadFixture("have some companies")
	company := &models.Company{}
	err := as.DB.First(company)
	as.NoError(err)
	as.NotZero(company.ID)
	as.NotZero(company.CreatedAt)

	company.Description = nulls.NewString("Now there are two.")

	res := as.HTML("/companies/%s", company.ID).Put(company)
	as.Equal(303, res.Code)
	as.Equal(fmt.Sprintf("/companies/%s", company.ID), res.Location())

	err = as.DB.Reload(company)
	as.NoError(err)
	as.Equal(nulls.NewString("Now there are two."), company.Description)
}

func (as *ActionSuite) Test_CompaniesResource_Destroy() {
	as.LoadFixture("have some companies")
	company := &models.Company{}
	err := as.DB.First(company)
	as.NoError(err)
	as.NotZero(company.ID)
	as.NotZero(company.CreatedAt)
	res := as.HTML("/companies/937af041-43ba-45d0-87a1-6bb173011996").Delete()
	as.Equal(303, res.Code)
	as.Equal("/companies", res.Location())
	count, err := as.DB.Count(&models.Companies{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_CompaniesResource_New() {
	res := as.HTML("/companies/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New Company")
}

func (as *ActionSuite) Test_CompaniesResource_Edit() {
	as.LoadFixture("have some companies")
	company := &models.Company{}
	err := as.DB.First(company)
	as.NoError(err)
	res := as.HTML("/companies/937af041-43ba-45d0-87a1-6bb173011996").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Description for company #1")
}
