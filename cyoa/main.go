package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	filename, templateName := parseFlags()
	arcs := getArcsFromFile(filename)
	tmpl := getHtmlTemplate(templateName)
	storyArcHandler := createStoryArcHandler(arcs, tmpl)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", storyArcHandler)
}

func parseFlags() (string, string) {
	filenamePtr := flag.String("file", "gopher.json", "Json file name containing stories. Default is `gopher.json`.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	templateNamePtr := flag.String("template", "template.html", "HTML file for rendering template. Default is `template.html`.")
	flag.StringVar(templateNamePtr, "t", *templateNamePtr, "Alias for --template.")
	flag.Parse()

	return strings.TrimSpace(*filenamePtr), strings.TrimSpace(*templateNamePtr)
}

func getArcsFromFile(filename string) map[string]Arc {
	arcs := make(map[string]Arc)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file %s: %s\n", filename, err)
	}
	err = json.Unmarshal(data, &arcs)
	if err != nil {
		log.Fatalf("Error parsing Json: %s\n", err)
	}
	return arcs
}

func getHtmlTemplate(templateName string) *template.Template {
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		log.Fatalf("Error parsing HTML template: %s\n", err)
	}
	return tmpl
}

func createStoryArcHandler(arcs map[string]Arc, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.RequestURI, "/")
		if key == "style.css" {
			// Only serve style.css file
			http.ServeFile(w, r, key)
			return
		}
		if key == "" {
			// Set intro as default
			key = "intro"
		}
		arc, ok := arcs[key]
		if !ok {
			http.NotFound(w, r)
			return
		}
		err := tmpl.Execute(w, arc)
		if err != nil {
			log.Fatalf("Error rendering HTML: %s\n", err)
		}
	}
}
