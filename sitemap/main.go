package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a parsed link from HTML
type Link struct {
	Href string
	Text string
}

// Sitemap XML structures
type Loc struct {
	Value string `xml:",chardata"`
}

type URL struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

func main() {
	urlFlag := flag.String("url", "", "URL to build sitemap for")
	depthFlag := flag.Int("depth", -1, "Maximum depth to crawl (-1 for unlimited)")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide a URL using the -url flag")
		os.Exit(1)
	}

	links := buildSitemap(*urlFlag, *depthFlag)

	sitemap := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0, len(links)),
	}

	for link := range links {
		sitemap.URLs = append(sitemap.URLs, URL{Loc: link})
	}

	output, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		fmt.Printf("Error generating XML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", output)
}

// buildSitemap crawls the website and returns all discovered URLs
func buildSitemap(rootURL string, maxDepth int) map[string]bool {
	visited := make(map[string]bool)

	baseURL, err := url.Parse(rootURL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	// Use BFS for proper depth handling
	type queueItem struct {
		url   string
		depth int
	}

	queue := []queueItem{{url: rootURL, depth: 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Skip if already visited
		if visited[current.url] {
			continue
		}

		// Skip if exceeding max depth
		if maxDepth >= 0 && current.depth > maxDepth {
			continue
		}

		visited[current.url] = true

		// Get links from current page
		links, err := getLinks(current.url)
		if err != nil {
			fmt.Printf("Error getting links from %s: %v\n", current.url, err)
			continue
		}

		// Process each link
		for _, link := range links {
			resolvedURL := resolveURL(baseURL, link.Href)
			if resolvedURL != "" && sameDomain(baseURL, resolvedURL) && !visited[resolvedURL] {
				queue = append(queue, queueItem{url: resolvedURL, depth: current.depth + 1})
			}
		}
	}

	return visited
}

// getLinks extracts all links from an HTML page
func getLinks(pageURL string) ([]Link, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseLinks(resp.Body)
}

// parseLinks parses HTML and extracts all anchor tags
func parseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []Link
	links = visitNode(doc, links)
	return links, nil
}

// visitNode recursively visits HTML nodes to find anchor tags
func visitNode(node *html.Node, links []Link) []Link {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, Link{
					Href: attr.Val,
					Text: getNodeText(node),
				})
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = visitNode(c, links)
	}

	return links
}

// getNodeText extracts text content from an HTML node
func getNodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getNodeText(c)
	}
	return strings.TrimSpace(text)
}

// resolveURL converts relative URLs to absolute URLs
func resolveURL(base *url.URL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}

	// Skipping non-HTTP links
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") ||
		strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "#") {
		return ""
	}

	parsedHref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	resolvedURL := base.ResolveReference(parsedHref)

	// Remove fragment
	resolvedURL.Fragment = ""

	return resolvedURL.String()
}

// sameDomain checks if two URLs belong to the same domain
func sameDomain(base *url.URL, targetURL string) bool {
	target, err := url.Parse(targetURL)
	if err != nil {
		return false
	}

	return strings.ToLower(base.Host) == strings.ToLower(target.Host)
}
