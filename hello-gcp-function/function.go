package function

import (
	"fmt"
	"net/http"
)

func HelloGCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "Hello, GCP!\n")
	fmt.Fprint(w, "This is a simple Go Cloud Function running successfully!\n")
	fmt.Fprintf(w, "Request method: %s\n", r.Method)
	fmt.Fprintf(w, "Request URL: %s\n", r.URL.Path)
}
