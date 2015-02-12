package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"

	cfg "github.com/jdkanani/smalldocs/config"
)

//
//  Project: Represents a single project
//
type Project struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Title       string        `json:"title" bson:"title"`
	Timestamp   time.Time
}

func ProjectInit(dbSession *mgo.Session, config *cfg.Config) error {
	// get a connection
	conn := dbSession.Copy()
	defer conn.Close()

	collection := conn.DB(config.Get("db.database")).C("projects")

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
