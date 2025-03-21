package linkparser

import (
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// Nested links will be extracted as separate link.
// Upon rendering through HTML, nested links are created as sibling.
// Try opening up test.html and you see that the nested link is a sibling.
func ParseLinks(node *html.Node) []Link {
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
		links = append(links, ParseLinks(next)...)
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
