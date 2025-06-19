package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"sync"

	"github.com/gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3003, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		
		var wg sync.WaitGroup

		idsToFetch := int(float64(numStories) * 1.5)
		if len(ids) > idsToFetch {
			ids = ids[:idsToFetch]
		}
		type result struct {
			index int
			item item
			err error
		}
		resultCh := make(chan result, len(ids))

		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		for i, id := range ids {
			wg.Add(1)
			go func(index int, id int) {
				defer wg.Done()
				hnItem, err := client.GetItem(id)
				if err != nil {
					resultCh <- result{index, item{}, err}
					return
				}
				resultCh <- result{index, parseHNItem(hnItem), nil}
			}(i, id)
		}
		go func() {
			wg.Wait()
			close(resultCh)
		}()

		orderedItems := make([]item, len(ids))
		receivedCount := 0
		
		for res := range resultCh {
			receivedCount++
			if res.err != nil {
				log.Printf("Error fetching item %d: %v", ids[res.index], res.err)
			}else {
				orderedItems[res.index] = res.item
			}
		}

		var stories []item
		for _, item := range orderedItems {
			if isStoryLink(item) {
				stories = append(stories, item)
				if len(stories) >= numStories {
					break
				}
			}
		}
		

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}