package actions

import (
	"github.com/gobuffalo/buffalo"
	"log"
)

type Breadcrumb struct {
	Name string
	Path string
}

// Push name and path to breadcrumb stack
func pushBreadcrumb(c buffalo.Context, name string) error {
	br := make([]Breadcrumb, 0)
	session := c.Session()
	bcs := session.Get("breadcrumbs")
	log.Println("bcs", bcs)
	if bcs != nil {
		br, _ = bcs.([]Breadcrumb)
	}
	path := c.Request().RequestURI
	br = append(br, Breadcrumb{Name: name, Path: path})
	log.Println("crumbs", br)

	// if we find current name and path earlier in breadcrumbs
	// shorten breadcrumbs to that part
	for i, b := range br {
		if b.Name == name && b.Path == path {
			br = br[:i+1]
			break
		}
	}
	session.Set("breadcrumbs", br)
	err := session.Save()
	if err != nil {
		return err
	}
	c.Set("crumbs", br)
	return nil
}

// Push name and path to breadcrumb stack
func setBreadcrumbs(c buffalo.Context) {
	br := make([]Breadcrumb, 0)
	session := c.Session()
	bcs := session.Get("breadcrumbs")
	if bcs != nil {
		br, _ = bcs.([]Breadcrumb)
	}
	c.Set("crumbs", br)
}

// Push name and path to breadcrumb stack
func newBreadcrumbs(c buffalo.Context) error {
	br := make([]Breadcrumb, 0)
	session := c.Session()
	session.Set("breadcrumbs", br)
	err := session.Save()
	if err != nil {
		return err
	}
	c.Set("crumbs", br)
	return nil
}
