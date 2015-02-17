package controllers

import (
	"net/http"

	"github.com/loconsolutions/goa"
)

//
// Document index
//
func DocumentIndex(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectByName(ctx.Params["pname"])
	if err != nil {
		return 404, err
	}
	return 200, ctx.Render("docs", map[string]interface{}{
		"project": project,
	})
}

//
// Edit document index
//
func EditIndex(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectById(ctx.Params["pname"])
	if err != nil {
		return 404, err
	}
	return 200, ctx.Render("docs-edit", map[string]interface{}{
		"project": project,
	})
}
