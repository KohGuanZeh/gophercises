package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	filename := parseFlags()
	node := parseHtmlFile(filename)
	links := parseLinks(node)
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

// Nested links will be extracted as separate link.
// Upon rendering through HTML, nested links are created as sibling.
// Try opening up test.html and you see that the nested link is a sibling.
func parseLinks(node *html.Node) []Link {
	if node.Type == html.ElementNode && node.Data == "a" {
		href := ""
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}
		var sb strings.Builder
		for next := node.FirstChild; next != nil; next = next.NextSibling {
			sb.WriteString(getAnchorText(next))
		}
		return []Link{{Href: href, Text: sb.String()}}
	}
	links := make([]Link, 0)
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		links = append(links, parseLinks(next)...)
	}
	return links
}

func getAnchorText(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}
	var sb strings.Builder
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		sb.WriteString(getAnchorText(next))
	}
	return sb.String()
}

func printLinks(links []Link) {
	fmt.Println("\nLinks found:", len(links))
	fmt.Println()
	for _, link := range links {
		fmt.Printf("Href: %s\n", link.Href)
		fmt.Printf("Text: %s\n\n", link.Text)
	}
}
