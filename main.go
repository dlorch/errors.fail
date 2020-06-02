package main

import (
	"fmt"
	"github.com/dlorch/errors.fail/session"
	"log"
	"net/http"
	"os"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	probe := session.ReadBoolOrDefault("http_probe", true)

	fmt.Fprintf(w, "<!DOCTYPE html>\n")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "<title>errors.fail - A free service that provides probing errors to your monitoring solutions</title>\n")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body>\n")
	fmt.Fprintf(w, "<h1>errors.fail</h1>\n")
	fmt.Fprintf(w, "<p>A free service that provides probing errors to your monitoring solutions.</p>\n")
	fmt.Fprintf(w, "<h2>HTTPS Probe</h2>\n")
	fmt.Fprintf(w, "<p><b>Endpoint:</b> curl -v <a href=\"https://probe.errors.fail/?%s\">https://probe.errors.fail/?%s</a>\n", session.SessionID, session.SessionID)
	if probe {
		fmt.Fprintf(w, "<p><b>Setting:</b> 200 OK</p>\n")
	} else {
		fmt.Fprintf(w, "<p><b>Setting:</b> 500 Internal Server Error</p>\n")
	}
	fmt.Fprintf(w, "<form action=\"/settings?%s\" method=\"POST\">\n", session.SessionID)
	fmt.Fprintf(w, "<input type=\"checkbox\" name=\"http_probe\"")
	if probe {
		fmt.Fprintf(w, " checked")
	}
	fmt.Fprintf(w, "/>Enabled\n")
	fmt.Fprintf(w, "<input type=\"submit\" value=\"Change Setting\"/>\n")
	fmt.Fprintf(w, "</form>\n")
	fmt.Fprintf(w, "<h2>ICMP Probe</h2>\n")
	fmt.Fprintf(w, "<p><b>No artificial packet loss:</b> ping probe.errors.fail</p>\n")
	fmt.Fprintf(w, "<p><b>50%% packet loss:</b> ping packetloss.errors.fail</p>\n")
	fmt.Fprintf(w, "<h2>Soon: Expired TLS/SSL Certificate</h2>\n")
	fmt.Fprintf(w, "<p><b>Endpoint:</b> curl -v <a href=\"https://expired.errors.fail\">https://expired.errors.fail</a></p>\n")
	fmt.Fprintf(w, "<h2>Contact</h2>\n")
	fmt.Fprintf(w, "<p>Created by <a href=\"https://dlorch.me/\">Daniel Lorch</a> in 2020</p>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
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
	}
}

func main() {
	http.HandleFunc("/settings", session.WithSession(changeSettings))
	http.HandleFunc("/", session.WithSession(mainPage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
