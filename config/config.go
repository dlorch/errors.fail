package config

import "os"

var (
	ProjectID = os.Getenv("PROJECT_ID")
)
