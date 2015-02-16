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

// get topic by id
func TopicById(id string) (*models.Topic, error) {
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	var topic = new(models.Topic)
	collection := session.DB(db).C("topics")
	if err := collection.FindId(bson.ObjectIdHex(id)).One(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

//
// Get topic
//
func GetTopic(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	topic, err := TopicById(ctx.Params["tid"])
	if err != nil {
		return 404, err
	}

	return 200, ctx.JSON(topic)
}

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

	fmt.Println(name)
	var topic *models.Topic = new(models.Topic)
	if err := collection.Find(bson.M{
		"name":    name,
		"project": bson.ObjectIdHex(ctx.Params["pid"]),
	}).One(topic); err == nil {
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

//
// Save topic
//
func SaveTopic(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	userTopic := new(models.Topic)
	if err := ctx.ReadJson(userTopic); err != nil {
		return 500, err
	}

	topic, err := TopicById(ctx.Params["tid"])
	if err != nil {
		return 404, err
	}

	topic.Title = utils.Title(userTopic.Title)
	topic.Name = utils.Slug(userTopic.Title)
	if topic.Name == "" {
		return 412, fmt.Errorf("Invalid title for topic!")
	}

	query := bson.M{"_id": topic.ID}
	change := bson.M{"$set": bson.M{
		"name":      topic.Name,
		"title":     topic.Title,
		"timestamp": time.Now().Unix(),
	},
	}

	// get mongo database
	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()
	collection := session.DB(db).C("topics")
	if err := collection.Update(query, change); err != nil {
		fmt.Println(err)
		return 500, err
	}

	return 200, ctx.JSON(topic)
}

//
// Delete topic by id
//
func DeleteTopic(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	topic, err := TopicById(ctx.Params["tid"])
	if err != nil {
		return 404, err
	}

	var db = context.Config.Get("db.database")

	// get mongodb session
	session := context.DBSession.Copy()
	defer session.Close()

	// remove project from collection
	collection := session.DB(db).C("topics")
	if err := collection.RemoveId(topic.ID); err != nil {
		return 500, err
	}

	return 200, ctx.JSON(topic)
}
