package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/throskam/kix/sess"
)

var ErrCSRFTokenMismatch = errors.New("CSRF mismatch")

var CSRFTokenKey = "csrf-token"

// GenerateCSRFToken generates a new CSRF token.
func GenerateCSRFToken() string {
	b := make([]byte, 16)

	// never returns an error.
	_, _ = rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}

// GetCSRFToken returns the CSRF token from the session.
func GetCSRFToken(ctx context.Context) string {
	return sess.MustGetSession(ctx).GetFirst(CSRFTokenKey)
}

// RefreshCSRFToken refreshes the CSRF token in the session.
func RefreshCSRFToken(ctx context.Context) {
	session := sess.MustGetSession(ctx)

	token := GenerateCSRFToken()

	session.Reset(CSRFTokenKey, token)
}

// GetFreshCSRFToken returns a fresh CSRF token and stores it in the session.
func GetFreshCSRFToken(ctx context.Context) string {
	RefreshCSRFToken(ctx)

	return GetCSRFToken(ctx)
}

// GetCSRFTokenOrCreate returns the CSRF token from the session or creates a new one.
func GetCSRFTokenOrCreate(ctx context.Context) string {
	token := GetCSRFToken(ctx)

	if token == "" {
		return GetFreshCSRFToken(ctx)
	}

	return token
}
