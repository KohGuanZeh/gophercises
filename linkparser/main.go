package main

import (
	"flag"
	"fmt"
	"gophercises/linkparser/linkparser"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	filename := parseFlags()
	node := parseHtmlFile(filename)
	links := linkparser.ParseLinks(node)
	printLinks(links)
}

func parseFlags() string {
	filenamePtr := flag.String("file", "test.html", "HTML file to extract links.")
	flag.StringVar(filenamePtr, "f", *filenamePtr, "Alias for --file.")
	flag.Parse()
	return strings.TrimSpace(*filenamePtr)
}

func parseHtmlFile(filename string) *html.Node {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error reading file %s: %s\n", filename, err)
	}
	node, err := html.Parse(f)
	f.Close()
	if err != nil {
		log.Fatalf("Error parsing HTML: %s\n", err)
	}
	return node
}

func printLinks(links []linkparser.Link) {
	fmt.Println("\nLinks found:", len(links))
	fmt.Println()
	for _, link := range links {
		fmt.Printf("Href: %s\n", link.Href)
		fmt.Printf("Text: %s\n\n", link.Text)
	}
}
