package main

import (
	c "github.com/jdkanani/smalldocs/config"

	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Node
type Node struct {
	Name   string `json:"name"`
	IsTree bool   `json:"isTree"`
	Path   string `json:"path"`
}

// Configuration
var config *c.Config

var spaces = regexp.MustCompile(" +")

// Template utility
var templatePath = "templates"
var templates = template.Must(template.ParseGlob("templates/*"))

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	var err error
	switch status {
	case http.StatusNotFound:
		err = templates.ExecuteTemplate(w, "error", "404 resource not found")
	case http.StatusInternalServerError:
		err = templates.ExecuteTemplate(w, "error", "500 Internal server error")
	}

	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}

// Utility
func FormatFilename(f string) string {
	var result = strings.TrimSuffix(f, filepath.Ext(f))
	result = strings.Replace(strings.ToLower(result), "_", " ", -1)
	return strings.Title(result)
}

// Index handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// Get tree
func TreeHandler(w http.ResponseWriter, r *http.Request) {
	// clean tree path
	var treePath = path.Clean(strings.TrimPrefix(r.URL.Path, "/tree/"))
	var basePath = filepath.Join(config.RootDirectory, treePath)

	var err error

	switch r.Method {
	case "GET":
		var result = make([]Node, 0)
		files, err := ioutil.ReadDir(basePath)
		if err != nil {
			ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		for _, f := range files {
			if p := strings.TrimPrefix(f.Name(), config.RootDirectory); p != "" {
				result = append(result, Node{
					IsTree: f.IsDir(),
					Path:   filepath.Join(treePath, f.Name()),
					Name:   FormatFilename(f.Name()),
				})
			}
		}
		if jsonData, err := json.Marshal(result); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}

	default:
		ErrorHandler(w, r, http.StatusNotFound)
	}

	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
}

// Show blob
func BlobHandler(w http.ResponseWriter, r *http.Request) {
	var blobPath = path.Clean(strings.TrimPrefix(r.URL.Path, "/blob/"))
	blobPath = filepath.Join(config.RootDirectory, blobPath)

	var err error
	switch r.Method {
	case "GET":
		if output, err := ioutil.ReadFile(blobPath); err == nil {
			var needHTML = false

			if types, ok := r.URL.Query()["type"]; ok {
				for _, t := range types {
					if t == "html" {
						needHTML = true
						break
					}
				}
			}

			if needHTML {
				output = blackfriday.MarkdownCommon(output)
			}

			w.Write(output)
		}

	case "POST":
		newName := r.FormValue("name")
		markdown := r.FormValue("markdown")

		file, err := os.Open(blobPath)
		if os.IsNotExist(err) {
			file, err = os.Create(blobPath)
		}

		if err == nil {
			defer file.Close()

			fileInfo, _ := file.Stat()
			ioutil.WriteFile(blobPath, []byte(markdown), fileInfo.Mode())

			// rename blob
			dir := filepath.Dir(blobPath)
			newName = strings.Trim(newName, " ")
			newName = strings.Replace(strings.ToLower(newName), "_", " ", -1)
			newName = spaces.ReplaceAllString(newName, "_")
			newPath := filepath.Join(dir, newName) + ".md"
			os.Rename(blobPath, newPath)

			result := &Node{
				IsTree: false,
				Path:   strings.TrimPrefix(newPath, config.RootDirectory),
				Name:   FormatFilename(newName),
			}

			if jsonData, err := json.Marshal(result); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
			}
		}

	case "DELETE":
		os.RemoveAll(blobPath)

	default:
		ErrorHandler(w, r, http.StatusNotFound)
	}

	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

func main() {
	// Load configuration
	cf, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	config = c.New(cf)

	// Create router
	mux := http.NewServeMux()
	// Static file handlers
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	// tree and blob handlers
	mux.HandleFunc("/blob/", BlobHandler)
	mux.HandleFunc("/tree/", TreeHandler)

	// Root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		IndexHandler(w, r)
	})

	// Start server
	fmt.Printf("Listening on %s ...", config.HostName)
	log.Fatal(http.ListenAndServe(config.HostName, mux))
}
