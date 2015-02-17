package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/jdkanani/goa"

	"github.com/loconsolutions/smalldocs/context"
	"github.com/loconsolutions/smalldocs/models"
	"github.com/loconsolutions/smalldocs/utils"

	"labix.org/v2/mgo/bson"
)

// get Project by id
func ProjectById(id string) (*models.Project, error) {
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
func ProjectByName(name string) (*models.Project, error) {
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
// Project page
//
func ProjectIndex(context *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
	return 301, nil
}

//
// Project name check
//
func CheckProject(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var data = make(map[string]string)
	ctx.ReadJson(&data)
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
	ctx.JSON(&map[string]string{
		"title": title,
		"name":  name,
	})
	return 200, nil
}

//
// Get all projects
//
func GetAllProjects(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	db := context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("projects")

	var projects []models.Project = make([]models.Project, 0)
	if err := collection.Find(nil).Sort("-timestamp").All(&projects); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(&projects)
}

//
// Get project
//
func GetProject(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	isId, err := regexp.MatchString(`[A-Fa-f0-9]{24}`, ctx.Params["pid"])
	if err != nil {
		return 500, err
	}

	var project *models.Project

	if isId {
		project, err = ProjectById(ctx.Params["pid"])
	} else {
		project, err = ProjectByName(ctx.Params["pid"])
	}

	if err != nil {
		return 404, err
	}

	return 200, ctx.JSON(project)
}

//
// Create new project
//
func PostProject(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("projects")

	var project = new(models.Project)
	if err := ctx.ReadJson(project); err != nil {
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

	return 200, ctx.JSON(project)
}

//
// Save project id
//
func SaveProject(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	userProject := new(models.Project)
	if err := ctx.ReadJson(userProject); err != nil {
		return 500, err
	}

	project, err := ProjectById(ctx.Params["pid"])
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

	return 200, ctx.JSON(project)
}

//
// Delete project by id
//
func DeleteProject(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectById(ctx.Params["pid"])
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

	return 200, ctx.JSON(project)
}

//
// Project settings
//
func ProjectSetting(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectByName(ctx.Params["pname"])
	if err != nil {
		return 404, err
	}

	return 200, ctx.Render("projectSettings", &project)
}
