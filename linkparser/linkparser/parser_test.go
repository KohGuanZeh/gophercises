package linkparser_test

import (
	"gophercises/linkparser/linkparser"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func isCorrectLink(link linkparser.Link, href string, text string) bool {
	return link.Href == href && link.Text == text
}

func TestParseLinkNormal(t *testing.T) {
	testHtml := `
<html>
	<body>
		<a href="/test-1">First</a>
		<a href="/test-2">Second</a>
		<a href="/test-3">Third</a>
	</body>
</html>
	`
	node, err := html.Parse(strings.NewReader(testHtml))
	if err != nil {
		t.Errorf("Error parsing HTML string: %s", err)
	}
	links := linkparser.ParseLinks(node)
	if len(links) != 3 ||
		!isCorrectLink(links[0], "/test-1", "First") ||
		!isCorrectLink(links[1], "/test-2", "Second") ||
		!isCorrectLink(links[2], "/test-3", "Third") {
		t.Errorf("Incorrect Links Gathered.")
	}
}

func TestParseLinkNestedElementsInAnchor(t *testing.T) {
	testHtml := `
<html>
	<body>
		<a href="/normal">Normal<span>Nested Element</span><div>Another Nested Element</div></a>
	</body>
</html>
	`
	node, err := html.Parse(strings.NewReader(testHtml))
	if err != nil {
		t.Errorf("Error parsing HTML string: %s", err)
	}
	links := linkparser.ParseLinks(node)
	if len(links) != 1 ||
		!isCorrectLink(links[0], "/normal", "NormalNested ElementAnother Nested Element") {
		t.Errorf("Incorrect Links gathered")
	}
}

func TestParseLinkNestedAnchor(t *testing.T) {
	testHtml := `
<html>
	<body>
		<a href="/normal">Normal<a href="/nested">Nested</a></a>
	</body>
</html>
	`
	node, err := html.Parse(strings.NewReader(testHtml))
	if err != nil {
		t.Errorf("Error parsing HTML string: %s", err)
	}
	links := linkparser.ParseLinks(node)
	if len(links) != 2 ||
		!isCorrectLink(links[0], "/normal", "Normal") ||
		!isCorrectLink(links[1], "/nested", "Nested") {
		t.Errorf("Incorrect Links gathered")
	}
}
