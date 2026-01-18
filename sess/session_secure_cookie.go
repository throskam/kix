package sess

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	ErrSecureCookieSessionReadFailure  = errors.New("failed to read secure cookie session")
	ErrSecureCookieSessionWriteFailure = errors.New("failed to write secure cookie session")
	ErrSecureCookieSessionEraseFailure = errors.New("failed to erase secure cookie session")
)

// SecureCookieSessionStore is a session store that uses secure cookies.
type SecureCookieSessionStore struct {
	store *sessions.CookieStore
}

// NewSecureCookieSessionStore creates a new secure cookie session store.
func NewSecureCookieSessionStore(secretKey ...[]byte) *SecureCookieSessionStore {
	store := sessions.NewCookieStore(secretKey...)

	store.Options.HttpOnly = true

	return &SecureCookieSessionStore{store}
}

// Read reads the session from the request.
func (s *SecureCookieSessionStore) Read(r *http.Request) (*Session, error) {
	sess, err := s.store.Get(r, SessionCookieKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSecureCookieSessionReadFailure, err)
	}

	session := NewSession()

	for k, v := range sess.Values {
		key, ok := k.(string)
		if !ok {
			continue
		}

		values, ok := v.([]string)
		if !ok {
			continue
		}

		for _, value := range values {
			session.Add(key, value)
		}
	}

	return session, nil
}

// Write writes the session to the response.
func (s *SecureCookieSessionStore) Write(r *http.Request, w http.ResponseWriter, session *Session) error {
	sess, _ := s.store.Get(r, SessionCookieKey)

	for k, v := range session.values {
		if len(v) == 0 {
			delete(sess.Values, k)
			continue
		}
		sess.Values[k] = v
	}

	err := sess.Save(r, w)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSecureCookieSessionWriteFailure, err)
	}

	return nil
}

// Erase erases the session from the response.
func (s *SecureCookieSessionStore) Erase(r *http.Request, w http.ResponseWriter) error {
	session, _ := s.store.Get(r, SessionCookieKey)

	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSecureCookieSessionEraseFailure, err)
	}

	return nil
}
