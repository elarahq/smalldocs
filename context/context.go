package context

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/jdkanani/smalldocs/config"
	"labix.org/v2/mgo"
)

// context for application
type Context struct {
	// templates
	templates map[string]*template.Template
	layout    []string

	// configuration
	Config *config.Config

	// mongodb session
	DBSession *mgo.Session
}

// render template
func (this *Context) RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	if name == "" {
		return fmt.Errorf("Template name must not be empty")
	}

	if this.templates == nil {
		this.templates = make(map[string]*template.Template)
	}

	if this.layout == nil {
		layout, err := filepath.Glob(filepath.Join(this.Config.Get("app.templates"), "layouts/*.tmpl"))
		if err != nil {
			return err
		}
		this.layout = layout
	}

	// Ensure the template exists in the map.
	tmpl, ok := this.templates[name]
	if !ok {
		files := make([]string, 0, len(this.layout)+1)
		files = append(files, this.layout...)
		files = append(files, filepath.Join(this.Config.Get("app.templates"), name)+".tmpl")
		tmpl = template.Must(template.ParseFiles(files...))
		// add to cache
		this.templates[name] = tmpl
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		return err
	}
	return nil
}

/**
 * Send json to response
 */
func (this *Context) JSON(w http.ResponseWriter, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonData)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	return nil
}

//
// Send json to response
//
func (this *Context) Text(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
	return nil
}

// ReadJson will parses the JSON-encoded data in the http
// Request object and stores the result in the value
// pointed to by v.
func ReadJson(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
