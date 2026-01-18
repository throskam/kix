package sess

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrCookieSessionReadFailure = errors.New("failed to read cookie session")

// CookieSessionStore is a session store that uses cookies.
type CookieSessionStore struct{}

// NewCookieSessionStore creates a new cookie session store.
func NewCookieSessionStore() *CookieSessionStore {
	return &CookieSessionStore{}
}

// Read reads the session from the request.
func (s *CookieSessionStore) Read(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(SessionCookieKey)
	if err != nil {
		return NewSession(), nil
	}

	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCookieSessionReadFailure, err)
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCookieSessionReadFailure, err)
	}

	session := NewSession()
	session.values = values

	return session, nil
}

// Write writes the session to the response.
func (s *CookieSessionStore) Write(r *http.Request, w http.ResponseWriter, session *Session) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(session.values.Encode()))

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieKey,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return nil
}

// Erase erases the session from the response.
func (s *CookieSessionStore) Erase(r *http.Request, w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:   SessionCookieKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return nil
}
