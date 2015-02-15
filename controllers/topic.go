package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jdkanani/goa"

	"github.com/jdkanani/smalldocs/context"
	"github.com/jdkanani/smalldocs/models"
	"github.com/jdkanani/smalldocs/utils"

	"labix.org/v2/mgo/bson"
)

//
// Get project topics
//
func GetTopics(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	projectId := ctx.Params["pid"]
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	// remove project from collection
	collection := session.DB(db).C("topics")

	var topics = make([]models.Topic, 0)
	if err := collection.Find(bson.M{"project": bson.ObjectIdHex(projectId)}).All(&topics); err != nil {
		return 500, err
	}
	return 200, ctx.JSON(topics)
}

//
// Topic name check
//
func CheckTopic(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
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
	collection := session.DB(db).C("topics")

	var topic *models.Topic = new(models.Topic)
	if err := collection.Find(bson.M{"name": name, "project": bson.ObjectIdHex(ctx.Params["pid"])}).One(topic); err == nil {
		if topic.ID.Hex() != id {
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
// Create new topic
//
func PostTopic(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectById(ctx.Params["pid"])
	if err != nil {
		return 404, err
	}

	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	collection := session.DB(db).C("topics")

	var topic = new(models.Topic)
	if err := ctx.ReadJson(topic); err != nil {
		return 500, err
	}

	topic.Title = utils.Title(topic.Title)
	topic.Name = utils.Slug(topic.Title)
	topic.Project = project.ID
	topic.Timestamp = time.Now().Unix()
	if topic.Name == "" {
		return 412, fmt.Errorf("Invalid title for topic!")
	}

	id := bson.NewObjectId()
	topic.ID = id
	if err := collection.Insert(topic); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(topic)
}
