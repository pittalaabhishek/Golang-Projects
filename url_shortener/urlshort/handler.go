package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.

// YAML is expected to be in the format:

//   - path: /some-path
//     url: https://www.some-url.com/demo

// The only errors that can be returned all related to having
// invalid YAML data.

// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []pathURL
	if len(yml) == 0 {
		return nil, fmt.Errorf("YAML data is empty")
	}
	err := yaml.Unmarshal(yml, &pathURLs)

	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		fmt.Printf("YAML content: %s\n", string(yml))
		return nil, err
	}

	fmt.Printf("Parsed YAML length: %d\n", len(pathURLs))
	fmt.Printf("Parsed YAML: %+v\n", pathURLs)

	pathsToUrls := make(map[string]string)
	for _, pu := range pathURLs {
		pathsToUrls[pu.Path] = pu.URL
	}
	fmt.Println("Paths to URLs Mapin Handler:", pathsToUrls)

	return MapHandler(pathsToUrls, fallback), nil
}