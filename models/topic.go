package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/jdkanani/smalldocs/context"
)

//
//  Topic: Represents a single topic
//
type Topic struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Title     string        `json:"title" bson:"title"`
	Project   bson.ObjectId `json:"project" bson:"project"`
	Timestamp int64         `json:"timestamp" bson:"timestamp"`
}

func TopicInit() error {
	// get a connection
	conn := context.DBSession.Copy()
	defer conn.Close()

	collection := conn.DB(context.Config.Get("db.database")).C("topics")

	// Unique key index
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	return collection.EnsureIndex(index)
}
