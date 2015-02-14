package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	ctx "github.com/jdkanani/smalldocs/context"
	"github.com/jdkanani/smalldocs/models"
	"github.com/jdkanani/smalldocs/utils"

	"labix.org/v2/mgo/bson"
)

//
// Project page
//
func ProjectIndex(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	return 301, nil
}

// get Project by id
func ProjectById(context *ctx.Context, id string) (*models.Project, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var project = new(models.Project)
	collection := session.DB(db).C("projects")
	if err := collection.FindId(bson.ObjectIdHex(id)).One(project); err != nil {
		return nil, err
	}
	return project, nil
}

// get Project by name
func ProjectByName(context *ctx.Context, name string) (*models.Project, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var project = new(models.Project)
	collection := session.DB(db).C("projects")
	if err := collection.Find(bson.M{"name": name}).One(project); err != nil {
		return nil, err
	}
	return project, nil
}

//
// Project name check
//
func CheckProject(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var data = make(map[string]string)
	context.ReadJson(r, &data)
	id, _ := data["id"]
	title, ok := data["title"]
	if !ok {
		return 412, fmt.Errorf("Title is required")
	}

	name := utils.Slug(title)
	collection := session.DB(db).C("projects")

	var project *models.Project = new(models.Project)
	if err := collection.Find(bson.M{"name": name}).One(project); err == nil {
		if project.ID.Hex() != id {
			return 403, nil
		}
	}

	// send name
	context.JSON(w, &map[string]string{
		"title": title,
		"name":  name,
	})
	return 200, nil
}

//
// Get all projects
//
func GetAllProjects(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("projects")

	var projects []models.Project = make([]models.Project, 0)
	if err := collection.Find(nil).Sort("-timestamp").All(&projects); err != nil {
		return 500, err
	}

	return 200, context.JSON(w, &projects)
}

//
// Get project by Id
//
func GetProject(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/projects/(?P<pid>`+ID+`)/?$`))
	project, err := ProjectById(context, params["pid"])
	if err != nil {
		return 404, err
	}

	return 200, context.JSON(w, project)
}

//
// Create new project
//
func PostProject(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("projects")

	var project = new(models.Project)
	if err := context.ReadJson(r, project); err != nil {
		return 500, err
	}

	project.Title = utils.Title(project.Title)
	project.Name = utils.Slug(project.Title)
	project.Timestamp = time.Now().Unix()
	if project.Name == "" {
		return 412, fmt.Errorf("Invalid title for project!")
	}

	id := bson.NewObjectId()
	project.ID = id
	if err := collection.Insert(project); err != nil {
		return 500, err
	}

	return 200, context.JSON(w, project)
}

//
// Save project id
//
func SaveProject(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	userProject := new(models.Project)
	if err := context.ReadJson(r, userProject); err != nil {
		return 500, err
	}

	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/projects/(?P<pid>`+ID+`)/?$`))
	project, err := ProjectById(context, params["pid"])
	if err != nil {
		return 404, err
	}

	project.Title = utils.Title(userProject.Title)
	project.Name = utils.Slug(userProject.Title)
	project.Description = userProject.Description
	if project.Name == "" {
		return 412, fmt.Errorf("Invalid title for project!")
	}

	query := bson.M{"_id": project.ID}
	change := bson.M{"$set": bson.M{
		"name":        project.Name,
		"title":       project.Title,
		"description": project.Description,
		"timestamp":   time.Now().Unix(),
	},
	}

	// get mongo database
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()
	collection := session.DB(db).C("projects")
	if err := collection.Update(query, change); err != nil {
		return 500, err
	}

	return 200, context.JSON(w, project)
}

//
// Delete project by id
//
func DeleteProject(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/projects/(?P<pid>`+ID+`)/?$`))
	project, err := ProjectById(context, params["pid"])
	if err != nil {
		return 404, err
	}

	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	// remove project from collection
	collection := session.DB(db).C("projects")
	if err := collection.RemoveId(project.ID); err != nil {
		return 500, err
	}

	return 200, context.JSON(w, project)
}

//
// Project settings
//
func ProjectSetting(context *ctx.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	params := utils.GetMatchedParams(r.URL.Path, regexp.MustCompile(`/projects/(?P<pname>`+SLUG+`)/settings/?$`))

	project, err := ProjectByName(context, params["pname"])
	if err != nil {
		return 404, err
	}

	return 200, context.RenderTemplate(w, "projectSettings", &project)
}
