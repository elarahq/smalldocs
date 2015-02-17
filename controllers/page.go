package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jdkanani/goa"

	"github.com/loconsolutions/smalldocs/context"
	"github.com/loconsolutions/smalldocs/models"
	"github.com/loconsolutions/smalldocs/utils"

	"labix.org/v2/mgo/bson"
)

// get page by id
func PageById(id string) (*models.Page, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var page = new(models.Page)
	collection := session.DB(db).C("pages")
	if err := collection.FindId(bson.ObjectIdHex(id)).One(page); err != nil {
		return nil, err
	}
	return page, nil
}

//
// Get page
//
func GetPage(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	page, err := PageById(ctx.Params["pageId"])
	if err != nil {
		return 404, err
	}

	return 200, ctx.JSON(page)
}

//
// Get pages
//
func GetPages(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	// remove project from collection
	collection := session.DB(db).C("pages")

	var pages = make([]models.Page, 0)
	if err := collection.Find(bson.M{
		"project": bson.ObjectIdHex(ctx.Params["pid"]),
		"topic":   bson.ObjectIdHex(ctx.Params["tid"]),
	}).All(&pages); err != nil {
		return 500, err
	}
	return 200, ctx.JSON(pages)
}

//
// Page name check
//
func CheckPage(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
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
	collection := session.DB(db).C("pages")

	var page *models.Page = new(models.Page)
	if err := collection.Find(bson.M{
		"name":    name,
		"project": bson.ObjectIdHex(ctx.Params["pid"]),
		"topic":   bson.ObjectIdHex(ctx.Params["tid"]),
	}).One(page); err == nil {
		if page.ID.Hex() != id {
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
// Create new page
//
func PostPage(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectById(ctx.Params["pid"])
	if err != nil {
		return 404, err
	}
	topic, err := TopicById(ctx.Params["tid"])
	if err != nil {
		return 404, err
	}

	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("pages")

	var page = new(models.Page)
	if err := ctx.ReadJson(page); err != nil {
		return 500, err
	}

	page.Title = utils.Title(page.Title)
	page.Name = utils.Slug(page.Title)
	page.Project = project.ID
	page.Topic = topic.ID
	page.Timestamp = time.Now().Unix()
	if page.Name == "" {
		return 412, fmt.Errorf("Invalid title for page!")
	}

	id := bson.NewObjectId()
	page.ID = id
	if err := collection.Insert(page); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(page)
}

//
// Save topic
//
func SavePage(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	userPage := new(models.Page)
	if err := ctx.ReadJson(userPage); err != nil {
		return 500, err
	}

	page, err := PageById(ctx.Params["pageId"])
	if err != nil {
		return 404, err
	}

	page.Title = utils.Title(userPage.Title)
	page.Name = utils.Slug(userPage.Title)
	page.Content = userPage.Content
	if page.Name == "" {
		return 412, fmt.Errorf("Invalid title for page!")
	}

	query := bson.M{"_id": page.ID}
	change := bson.M{"$set": bson.M{
		"name":      page.Name,
		"title":     page.Title,
		"content":   page.Content,
		"timestamp": time.Now().Unix(),
	},
	}

	// get mongo database
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()
	collection := session.DB(db).C("pages")
	if err := collection.Update(query, change); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(page)
}

//
// Delete page by id
//
func DeletePage(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	page, err := PageById(ctx.Params["pageId"])
	if err != nil {
		return 404, err
	}

	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	// remove project from collection
	collection := session.DB(db).C("pages")
	if err := collection.RemoveId(page.ID); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(page)
}
