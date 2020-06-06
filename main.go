package main

import (
	"fmt"
	"github.com/dlorch/errors.fail/session"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	pageTemplate string
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	values := struct {
		SessionID string
		HttpProbe bool
	}{
		SessionID: session.SessionID,
		HttpProbe: session.ReadBoolOrDefault("http_probe", true),
	}

	tmpl, err := template.New("pageTemplate").Parse(pageTemplate)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, values)
	if err != nil {
		log.Fatal(err)
	}
}

func changeSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Could not parse form values: %v", err)
			return
		}

		probe := r.FormValue("http_probe") != ""
		session.SaveBool("http_probe", probe, r)

		w.Header().Add("Location", fmt.Sprintf("/?%s", session.SessionID))
		w.WriteHeader(http.StatusFound)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
	}
}

func newSession(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		session.SessionID = session.GenerateUnsafeSessionID()

		w.Header().Add("Location", fmt.Sprintf("/?%s", session.SessionID))
		w.WriteHeader(http.StatusFound)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "400 Bad Request")
	}
}

func main() {
	http.HandleFunc("/settings", session.WithSession(changeSettings))
	http.HandleFunc("/new_session", session.WithSession(newSession))
	http.HandleFunc("/", session.WithSession(mainPage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	b, err := ioutil.ReadFile("template.html")
	if err != nil {
		log.Fatal(err)
	}
	pageTemplate = string(b)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
