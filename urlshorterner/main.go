package main

import (
	"flag"
	"fmt"
	"gophercises/urlshorterner/urlshort"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filename := parseFlags()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlData := getYamlData(filename)
	fmt.Printf("%s\n", yamlData)
	yamlHandler, err := urlshort.YAMLHandler(yamlData, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func parseFlags() string {
	filenamePtr := flag.String("file", "", "Yaml file name containing path and url pairs. Default is empty.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	flag.Parse()

	return strings.TrimSpace(*filenamePtr)
}

func getYamlData(filename string) []byte {
	if strings.TrimSpace(filename) == "" {
		return []byte(`
- path: /urlshort  
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file %s: %s\n", filename, err)
	}
	return data
}
