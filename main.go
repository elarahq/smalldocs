package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	cfg "github.com/jdkanani/smalldocs/config"
	ctx "github.com/jdkanani/smalldocs/context"

	"github.com/jdkanani/smalldocs/controllers"
	"github.com/jdkanani/smalldocs/handlers"
	"github.com/jdkanani/smalldocs/models"
	"github.com/jdkanani/smalldocs/router"
	"github.com/jdkanani/smalldocs/utils"

	"labix.org/v2/mgo"
)

func DBInit(config *cfg.Config) *mgo.Session {
	// Mongodb connection
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    strings.Split(config.Get("db.hosts"), ","),
		Timeout:  60 * time.Second,
		Database: config.Get("db.database"),
		// Username: Config.Get("db.username"),
		// Password: Config.Get("db.password"),
	}

	// create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	// Project init
	models.ProjectInit(mongoSession, config)

	return mongoSession
}

func main() {
	// get current directory
	root, err := os.Getwd()
	utils.Check(err)

	// load configuration
	config := new(cfg.Config)
	err = config.Load(filepath.Join(root, "config.ini"))
	utils.Check(err)

	// add root directory to config
	config.Set("app.root", root)
	config.Set("app.templates", filepath.Join(root, config.Get("app.templates")))
	config.Set("app.static", filepath.Join(root, config.Get("app.static")))

	// context
	context := &ctx.Context{
		Config:    config,
		DBSession: DBInit(config),
	}

	var AppHandlerFunc = func(fn handlers.HandleFunc) http.Handler {
		return &handlers.ErrorHandler{
			Context: context,
			Handler: fn,
		}
	}

	// static content handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// slug regexp
	slug := controllers.SLUG
	id := controllers.ID

	// router
	mux := new(router.Router)
	mux.NotFound(AppHandlerFunc(controllers.NotFound))

	mux.Get("/$", AppHandlerFunc(controllers.Index))

	// projects routes
	mux.Get("/projects/?$", AppHandlerFunc(controllers.ProjectIndex))
	mux.Get("/projects/all/?$", AppHandlerFunc(controllers.GetAllProjects))
	mux.Get("/projects/"+id+"/?$", AppHandlerFunc(controllers.GetProject))
	mux.Get("/projects/"+slug+"/settings/?$", AppHandlerFunc(controllers.ProjectSetting))
	mux.Post("/projects_check/?$", AppHandlerFunc(controllers.CheckProject))
	mux.Post("/projects/?$", AppHandlerFunc(controllers.PostProject))
	mux.Put("/projects/"+id+"/?$", AppHandlerFunc(controllers.SaveProject))
	mux.Delete("/projects/"+id+"/?$", AppHandlerFunc(controllers.DeleteProject))

	// docs routes
	mux.Get("/docs/"+slug+"/?$", AppHandlerFunc(controllers.DocumentIndex))

	// add router to http handle
	http.Handle("/", mux)

	// Start server
	serverURL := config.Get("server.host") + config.Get("server.port")
	fmt.Printf("Listening on %s ...", serverURL)
	log.Fatal(http.ListenAndServe(serverURL, nil))
}
