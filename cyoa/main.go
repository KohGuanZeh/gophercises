package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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
	filename := parseFlags()
	arcs := getArcsFromFile(filename)
	for key, value := range arcs {
		fmt.Printf("Arc: %s\n", key)
		fmt.Printf("Title: %s\n", value.Title)
		for _, story := range value.Story {
			fmt.Println(story)
		}
		for i, option := range value.Options {
			fmt.Printf("Option %d: %s (arc: %s)\n", i, option.Text, option.Arc)
		}
		fmt.Println()
	}
}

func parseFlags() string {
	filenamePtr := flag.String("file", "gopher.json", "Json file name containing stories. Default is `gopher.json`.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	flag.Parse()

	return strings.TrimSpace(*filenamePtr)
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
