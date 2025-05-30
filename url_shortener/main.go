package main
import (
    "fmt"
    "net/http"
	
	"url_shortener/urlshort"
)

func main() {
    mux := defaultMux()
    
    pathsToUrls := map[string]string{
        "/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
        "/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
    }
    mapHandler := urlshort.MapHandler(pathsToUrls, mux)
    
    // FIXED: Proper indentation with 2 spaces for 'url'
    yaml := `- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`
    
    fmt.Println("Paths to URLs Map:", pathsToUrls)
    fmt.Println("YAML Input:\n", yaml)
    
    yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Starting the server on :8080")
    http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, world!")
}

// yamlHandler (handles /urlshort, /urlshort-final) as we gave yamlHandler as the root handler to the server
//     ↓ (if no match)
// mapHandler (handles /urlshort-godoc, /yaml-godoc) 
//     ↓ (if no match)  
// defaultMux (handles everything else with "Hello, world!")