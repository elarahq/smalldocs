package controllers

import (
	"net/http"

	ctx "github.com/jdkanani/smalldocs/context"
	"github.com/jdkanani/smalldocs/models"
)

func ProjectIndex(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	return 301, nil
}

func GetProjects(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("projects")

	var projects []models.Project = make([]models.Project, 0)
	if err := collection.Find(nil).All(&projects); err != nil {
		return 500, err
	}

	return 200, context.JSON(w, &projects)
}
