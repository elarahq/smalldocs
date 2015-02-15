package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdkanani/smalldocs/config"
	"github.com/jdkanani/smalldocs/context"

	"github.com/jdkanani/smalldocs/controllers"
	"github.com/jdkanani/smalldocs/models"
	"github.com/jdkanani/smalldocs/utils"

	"github.com/jdkanani/goa"

	"labix.org/v2/mgo"
)

func DBInit() *mgo.Session {
	// Mongodb connection
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    strings.Split(context.Config.Get("db.hosts"), ","),
		Timeout:  60 * time.Second,
		Database: context.Config.Get("db.database"),
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
	models.ProjectInit(mongoSession)

	return mongoSession
}

func main() {
	// get current directory
	root, err := os.Getwd()
	utils.Check(err)

	// load configuration
	cfg := new(config.Config)
	context.Config = cfg
	err = context.Config.Load(filepath.Join(root, "config.ini"))
	utils.Check(err)

	// add root directory to config
	cfg.Set("app.root", root)
	cfg.Set("app.templates", filepath.Join(root, cfg.Get("app.templates")))
	cfg.Set("app.static", filepath.Join(root, cfg.Get("app.static")))

	context.DBSession = DBInit()

	// static content handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// router
	mux := new(goa.Router)
	mux.SetRenderer(&goa.Renderer{
		TemplateDir: context.Config.Get("app.templates"),
		LayoutDir:   context.Config.Get("app.layouts"),
	})
	mux.NotFound(controllers.NotFound)

	mux.Get("/", controllers.Index)

	// projects routes
	mux.Get("/projects/?$", controllers.ProjectIndex)
	mux.Get("/projects/all/?$", controllers.GetAllProjects)
	mux.Get("/projects/:pid/?$", controllers.GetProject)
	mux.Get("/projects/:pname/settings/?$", controllers.ProjectSetting)
	mux.Post("/projects_check/?$", controllers.CheckProject)
	mux.Post("/projects/?$", controllers.PostProject)
	mux.Put("/projects/:pid/?$", controllers.SaveProject)
	mux.Delete("/projects/:pid/?$", controllers.DeleteProject)

	// docs routes
	mux.Get("/docs/:pname/?$", controllers.DocumentIndex)
	mux.Get("/docs/:pname/:dname/:pagename/?$", controllers.PageIndex)

	// add router to http handle
	http.Handle("/", mux)

	// Start server
	serverURL := cfg.Get("server.host") + cfg.Get("server.port")
	fmt.Printf("Listening on %s ...", serverURL)
	log.Fatal(http.ListenAndServe(serverURL, nil))
}
