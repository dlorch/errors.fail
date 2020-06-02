package config

import "os"

var (
	ProjectID    = os.Getenv("PROJECT_ID")
	CookieDomain = os.Getenv("COOKIE_DOMAIN")
)
