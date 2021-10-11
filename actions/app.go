package actions

import (
	"github.com/gobuffalo/buffalo"
	//"github.com/gobuffalo/buffalo-goth/goth"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/unrolled/secure"

	"doors/models"

	"encoding/gob"

	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"

	"github.com/markbates/goth/gothic"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_doors_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())
		app.Use(Authorize())
		app.Use(SetCurrentUser())
		app.Middleware.Skip(Authorize(), HomeHandler)
		app.GET("/", HomeHandler)

		app.Resource("/doors", DoorsResource{})
		app.Resource("/companies", CompaniesResource{})
		p := app.Resource("/people", PeopleResource{})
		p.Resource("/tokens", TokensResource{})
		app.Resource("/access_groups", AccessGroupsResource{})

		agdr := &AccessGroupDoorsResource{}
		agd := app.Group("/access_group_doors")
		agd.POST("/", agdr.Create) // POST /access_group_doors => agdr.Create
		tagr := &TokenAccessGroupsResource{}
		tag := app.Group("/token_access_groups")
		tag.POST("/", tagr.Create) // POST /token_access_groups => tagr.Create
		auth := app.Group("/auth")
		bah := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("", AuthDestroy())
		auth.Middleware.Skip(Authorize(), bah, AuthCallback)
		auth.GET("/{provider}", bah)
		auth.GET("/{provider}/callback", AuthCallback)
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	// need to register []Breadcrumb{} to avoid
	// securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: []actions.Breadcrumb
	gob.Register([]Breadcrumb{})
	gob.Register(make(map[string]interface{}))
	gob.Register(models.User{})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
