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
// Document page
//
func DocumentIndex(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/docs/(?P<pname>`+SLUG+`)/?$`))
	fmt.Fprint(w, params["pname"])
	return 200, nil
}
