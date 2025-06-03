package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"adventure_story/templates"
)

type StoryHandler struct {
	story    Story
	template *template.Template
}

type Story map[string]Chapter

type Chapter struct {
	Title 	string 		`json:"title"`
	Story 	[]string 		`json:"story"`
	Options []Option   `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func NewStoryHandler(story Story, tmpl *template.Template) *StoryHandler {
	return &StoryHandler{
		story: story,
		template: tmpl,
	}
}

func main() {
	story, err := loadStoryFromFile("gopher.json")

	if err != nil {
		log.Fatalf("Error loading the file %v", err)
	}

	tmpl, err := template.New("story").Parse(templates.StoryTemplate)

	if err != nil {
		log.Fatalf("Error loading the html template %v", err)
	}

	storyHandler := NewStoryHandler(story, tmpl)

	http.Handle("/", storyHandler)

	fmt.Println("Starting Choose Your Own Adventure server on :8080")
	fmt.Println("Visit http://localhost:8080 to begin your adventure!")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadStoryFromFile(filename string) (Story, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open story file: %w", err)
	}
	defer file.Close()

	var story Story
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&story)

	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return story, nil
}

func (sh *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = strings.TrimPrefix(path, "/")

	chapter, exists := sh.story[path]
	if !exists {
		http.Error(w, "Chapter not found", http.StatusNotFound)
		return
	}

	err := sh.template.Execute(w, chapter)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
