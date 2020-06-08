package main

import (
	"encoding/json"
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

type settings struct {
	HTTP_Probe bool `json:"http_probe"`
}

func changeSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var s settings
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Fprintf(w, "Could not parse JSON values: %v", err)
			return
		}

		session.SaveBool("http_probe", s.HTTP_Probe, r)

		fmt.Fprintf(w, "OK")
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
