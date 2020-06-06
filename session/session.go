package session

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/dlorch/errors.fail/config"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	sessionLength = 6
)

var (
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx                   = context.Background()
	client                = createFirestoreClient(ctx)
	SessionID  string
)

func createFirestoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func WithSession(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// try to retrieve session ID from URL query (part after ? in URL)
		SessionID = r.URL.RawQuery
		if !isValidSessionID(SessionID) {
			// try to retrieve session ID from cookie
			sessionCookie, err := r.Cookie("sessionID")
			if err == nil {
				SessionID = sessionCookie.Value
			}

			// no valid session ID found -> create new session
			if !isValidSessionID(SessionID) {
				SessionID = GenerateUnsafeSessionID()
			}

			w.Header().Add("Location", fmt.Sprintf("%s?%s", r.URL.Path, SessionID))
			w.WriteHeader(http.StatusFound)
			return
		}

		sessionCookie := http.Cookie{
			Name:     "sessionID",
			Value:    SessionID,
			Domain:   config.CookieDomain, // "if a domain is specified, then subdomains are always included" (Set-Cookie - HTTP | MDN)
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &sessionCookie)

		handler(w, r)
	}
}

func GenerateUnsafeSessionID() string {
	b := make([]byte, sessionLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func isValidSessionID(session string) bool {
	if len(session) != sessionLength {
		return false
	}

	for i := 0; i < len(session); i++ {
		if strings.Index(charset, string(session[i])) == -1 {
			return false
		}
	}

	return true
}

func SaveBool(name string, value bool, r *http.Request) {
	_, err := client.Collection("sessions").Doc(SessionID).Set(ctx, map[string]interface{}{
		name:          value,
		"last_active": time.Now(),
		"user_agent":  r.UserAgent(),
		"remote_addr": r.RemoteAddr,
	}, firestore.MergeAll)
	if err != nil {
		log.Fatalf("Failed setting value for session \"%s\": %v", SessionID, err)
	}
}

func ReadBoolOrDefault(name string, defaultValue bool) bool {
	dsnap, err := client.Collection("sessions").Doc(SessionID).Get(ctx)
	if err != nil {
		return defaultValue
	}
	result := dsnap.Data()["http_probe"].(bool)
	return result
}
