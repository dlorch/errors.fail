package session

import (
	"testing"
)

func TestValidSessionID(t *testing.T) {
	validIDs := []string{
		"x0M10V", "O8yTPp", "DBr1a2", "H2Erq2", "DnMlUQ",
		"xh66Yh", "uccHJC", "ZvcsNF", "V2dP61", "S3ZfIh",
	}
	for _, sessionID := range validIDs {
		if !isValidSessionID(sessionID) {
			t.Errorf("Session ID \"%s\" is valid, but got invalid", sessionID)
		}
	}
}

func TestInvalidSessionIDs(t *testing.T) {
	invalidIDs := []string{
		"", "abc", "x0M-0V", "DB_1a2",
	}
	for _, sessionID := range invalidIDs {
		if isValidSessionID(sessionID) {
			t.Errorf("Session ID \"%s\" is invalid, but got valid", sessionID)
		}
	}
}
