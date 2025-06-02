package main

import (
	"fmt"
	"log"
	"os"

	"html_link_parser/parser"
)

func main() {
	htmlFile, err := os.Open("examples/ex2.html")

	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}
	defer htmlFile.Close()

	links, err := parser.ParseHtml(htmlFile)

	if err != nil {
		log.Fatalf("Failed to create links %v", err)
	}
	for _, link := range links {
		fmt.Printf("Href: %s Text: %s\n", link.Href, link.Text)
	}
}