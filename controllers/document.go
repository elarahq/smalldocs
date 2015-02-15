package controllers

import (
	"fmt"
	"net/http"
	"regexp"

	ctx "github.com/jdkanani/smalldocs/context"
	_ "github.com/jdkanani/smalldocs/models"
	"github.com/jdkanani/smalldocs/utils"

	_ "labix.org/v2/mgo/bson"
)

//
// Document index
//
func DocumentIndex(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/docs/(?P<pname>`+SLUG+`)/?$`))
	project, err := ProjectByName(context, params["pname"])
	if err != nil {
		return 404, err
	}
	return 200, context.RenderTemplate(w, "docs", map[string]interface{}{
		"project": project,
	})
}

//
//  Page index
//
func PageIndex(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	urlRegexp := fmt.Sprintf(`/docs/(?P<projectName>%s)/(?P<docName>%s)/(?P<pageName>%s)`, SLUG, SLUG, SLUG)
	_ = utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(urlRegexp))
	return 200, nil
}
