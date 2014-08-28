package main

import (
	"flag"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates *template.Template
var templateDir string
var staticDir string

func loadTemplates(r *mux.Router) error {
	log.Println("Loading templates:")
	templates = template.New("templates")
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		t := template.Must(template.ParseFiles(path))
		log.Printf("  %v\n", rel)
		templates = template.Must(templates.AddParseTree(rel, t.Tree))
		return nil
	})
	if err != nil {
		log.Fatal("Unable to load templates.")
	}

	log.Println("Loading static files:")
	err = filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}
		log.Printf("  %v\n", rel)

		r.HandleFunc("/"+rel, func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				writeJsonError(w, err)
				return
			}
			w.Write(b)
		}).Methods("GET")

		return nil
	})
	return err
}

func templateFlags() {
	flag.StringVar(&templateDir, "templates", "./templates", "The directory to search for template files.")
	flag.StringVar(&staticDir, "static", "./static", "The directory to search for static files.")
}
