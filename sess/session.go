package sess

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
)

// SessionCookieKey is the name of the session cookie.
var SessionCookieKey = "session"

var (
	ErrSessionEraseFailure = errors.New("failed to erase session")
	ErrSessionReadFailure  = errors.New("failed to read session")
	ErrSessionWriteFailure = errors.New("failed to write session")
)

// SessionStore is the interface for session stores.
type SessionStore interface {
	Read(*http.Request) (*Session, error)
	Write(*http.Request, http.ResponseWriter, *Session) error
	Erase(*http.Request, http.ResponseWriter) error
}

// Session is a session.
type Session struct {
	values url.Values
}

// NewSession creates a new session.
func NewSession() *Session {
	return &Session{
		values: url.Values{},
	}
}

// Get returns the values of the session for the given key.
func (s *Session) Get(key string) []string {
	return s.values[key]
}

// GetFirst returns the first value of the session for the given key.
func (s *Session) GetFirst(key string) string {
	return s.values.Get(key)
}

// Empty returns true if the session does not contain the given key.
func (s *Session) Empty(key string) bool {
	return !s.values.Has(key)
}

// Has returns true if the session contains the given key and value.
func (s *Session) Has(key, value string) bool {
	if s.Empty(key) {
		return false
	}

	return slices.Contains(s.values[key], value)
}

// Set sets the values of the session for the given key.
func (s *Session) Set(key string, values []string) {
	s.values[key] = values
}

// Reset sets the value of the session for the given key.
func (s *Session) Reset(key, value string) {
	s.values.Set(key, value)
}

// Add adds the value to the session for the given key.
func (s *Session) Add(key, value string) {
	if s.Has(key, value) {
		return
	}

	s.values.Add(key, value)
}

// Remove removes the value from the session for the given key.
func (s *Session) Remove(key, value string) {
	if !s.Has(key, value) {
		return
	}

	values := []string{}

	for _, v := range s.values[key] {
		if v != value {
			values = append(values, v)
		}
	}

	s.Set(key, values)
}

// Toggle toggles the value of the session for the given key.
func (s *Session) Toggle(key, value string) {
	if s.Has(key, value) {
		s.Remove(key, value)
	} else {
		s.Add(key, value)
	}
}

// Del deletes the session for the given key.
func (s *Session) Del(key string) {
	s.values.Del(key)
}
