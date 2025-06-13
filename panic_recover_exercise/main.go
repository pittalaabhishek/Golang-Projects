package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func recoverWrap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stackTrace := debug.Stack()

				log.Printf("PANIC: %v\n%s", rec, stackTrace)

				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusInternalServerError)

				isDev := os.Getenv("ENV") == "development"

				if isDev {
					fmt.Fprintf(w, "Error: %v\n\nStack Trace:\n%s", rec, stackTrace)
				} else {
					fmt.Fprint(w, "Something went wrong")
				}
			}
		}()

		handler(w, r)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", recoverWrap(panicDemo))
	mux.HandleFunc("/panic-after/", recoverWrap(panicAfterDemo))
	mux.HandleFunc("/", recoverWrap(hello))

	log.Fatal(http.ListenAndServe(":3006", mux))
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
