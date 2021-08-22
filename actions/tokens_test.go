package actions

import (
	"doors/models"

	"github.com/gofrs/uuid"
)

// TODO:
// ajax put => index / show
// unique constraint token_id

func (as *ActionSuite) Test_TokensResource_List() {
	as.LoadFixture("have some people")
	as.LoadFixture("have some tokens")
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/").Get()

	as.Equal(302, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())
}

func (as *ActionSuite) Test_TokensResource_Show() {
	as.LoadFixture("have some people")
	as.LoadFixture("have some tokens")
	token := &models.Token{}
	err := as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/3eb3735a-4ac1-44fc-9154-46fe073fb394").Get()
	as.Equal(302, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())
	res = as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Token#1")
	as.Contains(body, "Token#2")
}

func (as *ActionSuite) Test_TokensResource_Create() {
	as.LoadFixture("have some people")
	id, _ := uuid.NewV4()
	pID, _ := uuid.FromString("bd42798a-77cb-440c-9595-ec166fd3c32d")
	token := &models.Token{
		ID:       id,
		PersonID: pID,
		TokenID:  "abcdef",
	}
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/").Post(token)
	as.Equal(303, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())

	err := as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	as.NotZero(token.PersonID)
}

func (as *ActionSuite) Test_TokensResource_Create_Errors() {
	as.LoadFixture("have some people")
	token := &models.Token{
		ID:       uuid.UUID{},
		PersonID: uuid.UUID{},
	}
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/").Post(token)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "TokenID can not be blank.")

	c, err := as.DB.Count(token)
	as.NoError(err)
	as.Equal(0, c)
}

// Create duplicate Token with same TokenID leads to redirect_edit

func (as *ActionSuite) Test_TokensResource_Create_Duplicates() {
	as.LoadFixture("have some people")
	id, _ := uuid.NewV4()
	pID, _ := uuid.FromString("bd42798a-77cb-440c-9595-ec166fd3c32d")
	token := &models.Token{
		ID:       id,
		PersonID: pID,
		TokenID:  "abcdef",
	}
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/").Post(token)
	as.Equal(303, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())

	err := as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	as.NotZero(token.PersonID)

	id, _ = uuid.NewV4()
	token = &models.Token{
		ID:       id,
		PersonID: pID,
		TokenID:  "abcdef",
	}
	res = as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/").Post(token)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "TokenID \"abcdef\" is already in use")
}

func (as *ActionSuite) Test_TokensResource_Update() {
	as.LoadFixture("have some people")
	as.LoadFixture("have some tokens")
	token := &models.Token{}
	err := as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	as.NotZero(token.CreatedAt)

	token.TokenID = "99999999"

	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/%s", token.ID).Put(token)
	as.Equal(303, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())

	err = as.DB.Reload(token)
	as.NoError(err)
	as.Equal("99999999", token.TokenID)
}

func (as *ActionSuite) Test_TokensResource_Destroy() {
	as.LoadFixture("have some people")
	as.LoadFixture("have some tokens")
	token := &models.Token{}
	err := as.DB.First(token)
	as.NoError(err)
	as.NotZero(token.ID)
	as.NotZero(token.CreatedAt)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/3eb3735a-4ac1-44fc-9154-46fe073fb394/").Delete()
	as.Equal(303, res.Code)
	as.Equal("/people/bd42798a-77cb-440c-9595-ec166fd3c32d", res.Location())
	count, err := as.DB.Count(&models.Tokens{})
	as.NoError(err)
	as.Equal(3, count)
}

func (as *ActionSuite) Test_TokensResource_New() {
	as.LoadFixture("have some people")
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d/tokens/new").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "New Token")
}

func (as *ActionSuite) Test_TokensResource_Edit() {
	as.LoadFixture("have some people")
	as.LoadFixture("have some tokens")
	token := &models.Token{}
	err := as.DB.First(token)
	as.NoError(err)
	res := as.HTML("/people/bd42798a-77cb-440c-9595-ec166fd3c32d").Get()
	as.Equal(200, res.Code)
	body := res.Body.String()
	as.Contains(body, "Token#1")
}
