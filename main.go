package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	cfg "github.com/jdkanani/smalldocs/config"
	ctx "github.com/jdkanani/smalldocs/context"

	"github.com/jdkanani/smalldocs/controllers"
	"github.com/jdkanani/smalldocs/handlers"
	"github.com/jdkanani/smalldocs/router"
	"github.com/jdkanani/smalldocs/utils"
)

// App configuration
var Config *cfg.Config

func main() {
	// get current directory
	root, err := os.Getwd()
	utils.Check(err)

	// load configuration
	Config = new(cfg.Config)
	err = Config.Load(filepath.Join(root, "config.ini"))
	utils.Check(err)

	// add root directory to config
	Config.Set("app.root", root)
	Config.Set("app.templates", filepath.Join(root, Config.Get("app.templates")))
	Config.Set("app.static", filepath.Join(root, Config.Get("app.static")))

	// context
	context := new(ctx.Context)
	context.Config = Config
	var AppHandlerFunc = func(fn handlers.HandleFunc) http.Handler {
		return &handlers.ErrorHandler{
			Context: context,
			Handler: fn,
		}
	}

	// static content handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// router
	mux := new(router.Router)
	mux.Get("/", AppHandlerFunc(controllers.Index))

	// add router to http handle
	http.Handle("/", mux)

	// Start server
	serverURL := Config.Get("server.host") + Config.Get("server.port")
	fmt.Printf("Listening on %s ...", serverURL)
	log.Fatal(http.ListenAndServe(serverURL, nil))
}
