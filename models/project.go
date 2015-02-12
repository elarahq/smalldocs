package models

//
//  Project: Represents a single project
//
type Project struct {
	Name  string `json:"name" bson:"name"`
	Title string `json:"title" bson:"title"`
}
