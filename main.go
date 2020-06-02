package main

import (
	"fmt"
	"github.com/dlorch/errors.fail/session"
	"html"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	p, err := strconv.ParseFloat(r.URL.Path[1:], 32)
	if err != nil || p < 0 || p > 1.0 {
		p = 1.0
	}
	if rand.Float32() < float32(p) {
		fmt.Fprintf(w, r.URL.RawQuery)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>500 Internal Server Error</h1><p>Intentionally failing this request with probabilty of %.2f.</p>", p)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", session.WithSession(handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
