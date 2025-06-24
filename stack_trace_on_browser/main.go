package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func recoverWrap(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stackTrace := debug.Stack()
				log.Printf("PANIC: %v\n%s", rec, stackTrace)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Panic: %v\n\n", rec)
				fmt.Fprint(w, formatStackTrace(string(stackTrace)))
			}
		}()

		handler.ServeHTTP(w, r)
	}
}

func formatStackTrace(stack string) string {
	lines := strings.Split(stack, "\n")
	var formattedOutput strings.Builder

	re := regexp.MustCompile(`^(\s*)(\S+\.go):(\d+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			indent := matches[1]
			filePath := matches[2]
			lineNumber := matches[3]

			if _, err := os.Stat(filePath); err == nil {
				encodedFilePath := url.QueryEscape(filePath)
				debugLink := fmt.Sprintf("/debug/source?filepath=%s&line=%s", encodedFilePath, lineNumber)
				formattedOutput.WriteString(fmt.Sprintf("%s<a href=\"%s\">%s:%s</a>\n", indent, debugLink, filePath, lineNumber))
			} else {
				formattedOutput.WriteString(line + "\n")
			}
		} else {
			formattedOutput.WriteString(line + "\n")
		}
	}
	return formattedOutput.String()
}

func sourceHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("filepath")
	lineStr := r.URL.Query().Get("line")
	highlightLine := 0

	if line, err := strconv.Atoi(lineStr); err == nil {
		highlightLine = line
	}

	if filePath == "" {
		http.Error(w, "Missing 'filepath' parameter", http.StatusBadRequest)
		return
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	wd, _ := os.Getwd()
	if !strings.HasPrefix(absFilePath, wd) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	code, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusNotFound)
		return
	}

	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	
	var formatter *html.Formatter
	if highlightLine > 0 {
		formatter = html.New(
			html.WithLineNumbers(true),
			html.LineNumbersInTable(true),
			html.HighlightLines([][2]int{{highlightLine, highlightLine}}),
		)
	} else {
		formatter = html.New(
			html.WithLineNumbers(true),
			html.LineNumbersInTable(true),
		)
	}
	
	style := styles.GitHub

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error tokenizing code: %v", err), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error formatting code: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, buf.String())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", recoverWrap(panicDemo))
	mux.HandleFunc("/panic-after/", recoverWrap(panicAfterDemo))
	mux.HandleFunc("/", recoverWrap(hello))
	mux.HandleFunc("/debug/source", sourceHandler)

	log.Println("Server listening on :3006 (Set ENV=development for debug features)")
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
	panic("Oops! A controlled panic.")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello, world!</h1>")
}