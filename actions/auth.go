package actions

import (
	"doory/models"
	"fmt"
	"net/http"
	"strings"

	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/markbates/goth/providers/azureadv2"
)

func init() {
	var host string

	switch appHost := App().Host; strings.ContainsAny(appHost, "127.0.0.1") {
	case true:
		splits := strings.Split(appHost, ":")
		fmt.Println(splits)
		host = fmt.Sprintf("%s://localhost:%s", splits[0], splits[2])
	default:
		host = appHost
	}

	gothic.Store = App().SessionStore
	goth.UseProviders(
		azureadv2.New(
			os.Getenv("AZURE_KEY"),
			os.Getenv("AZURE_SECRET"),
			fmt.Sprintf(
				"%s%s",
				host,
				"/auth/azureadv2/callback",
			),
			azureadv2.ProviderOptions{
				Scopes: []azureadv2.ScopeType{
					azureadv2.AgreementReadAllScope,
				},
				Tenant: azureadv2.OrganizationsTenant,
			},
		),
	)
}

// AuthCallback handler for OAuth2 procedure
func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(
			fmt.Errorf("no transaction found"),
		)
	}
	exists, err := tx.Where("email = ?", user.Email).Exists(&models.User{})
	if err != nil {
		return errors.WithStack(err)
	}
	modelUser := &models.User{
		Name:       user.Name,
		Email:      nulls.NewString(user.Email),
		Provider:   user.Provider,
		ProviderID: user.UserID,
	}
	if exists {
		if err := tx.First(modelUser); err != nil {
			return errors.WithStack(err)
		}
	} else {
		if err := tx.Save(modelUser); err != nil {
			return errors.WithStack(
				fmt.Errorf(
					"could not create new user %s: %s",
					user.Name,
					err.Error(),
				),
			)
		}
	}
	c.Session().Set("current_user_id", modelUser.ID)
	if c.Session().Save() != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "Login successful!")
	return c.Redirect(
		http.StatusFound,
		"/",
	)
}

// SetCurrentUser sets the current user into the context
func SetCurrentUser() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			tx, ok := c.Value("tx").(*pop.Connection)
			if !ok {
				return errors.WithStack(
					fmt.Errorf("no transaction found"),
				)
			}
			sessionUserID := c.Session().Get("current_user_id")
			if sessionUserID == nil {
				return next(c)
			}
			user := &models.User{}
			if err := tx.Find(user, sessionUserID); err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", user)
			return next(c)
		}
	}
}

// Authorize blocks non-authorized endpoints
func Authorize() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			if c.Session().Get("current_user_id") == nil {
				return c.Redirect(
					http.StatusTemporaryRedirect,
					"/",
				)
			}
			return next(c)
		}
	}
}

// AuthDestroy logs out the user
func AuthDestroy() buffalo.Handler {
	return func(c buffalo.Context) error {
		if err := gothic.Logout(c.Response(), c.Request()); err != nil {
			c.Flash().Add("error", "Could not logout")
			c.Logger().Error("Could not write logout response:", err.Error())
			return c.Redirect(
				http.StatusInternalServerError,
				"/",
			)
		}
		c.Session().Clear()
		if err := c.Session().Save(); err != nil {
			c.Flash().Add("error", "Could not logout")
			c.Logger().Error("Could not write logout response:", err.Error())
			return errors.WithStack(err)
		}
		c.Flash().Add("success", "Logout successful")
		return c.Redirect(
			http.StatusFound,
			"/",
		)
	}
}
