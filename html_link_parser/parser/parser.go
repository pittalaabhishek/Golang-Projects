package parser

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseHtml(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)

	if err != nil {
		log.Fatalf("Could not parse the document")
	}

	var links []Link
	var extractLinks func(*html.Node)
	extractLinks = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			href := ""
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}
			if href != "" {
				LinkText := extractText(node)
				links = append(links, Link{Href: href, Text: LinkText})
			}
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extractLinks(child)
		}
	}
	extractLinks(doc)

	return links, nil
}

func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data)
	}

	var result strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "a" {
			continue
		}
		
		text := extractText(child)
		if text != "" {
			if result.Len() > 0 {
				result.WriteString(" ")
			}
			result.WriteString(text)
		}
	}
	return strings.TrimSpace(result.String())
}
