package sess

import (
	"context"
	"errors"
)

// contextKey is the type of the context key.
type contextKey string

// sessionContextKey is the context key for the session.
const sessionContextKey contextKey = "session"

var ErrCtxSessionMissing = errors.New("missing session context")

// GetSession returns the session from the context.
func GetSession(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(sessionContextKey).(*Session)
	if !ok {
		return nil, ErrCtxSessionMissing
	}

	return session, nil
}

// MustGetSession returns the session from the context or panics if it is missing.
func MustGetSession(ctx context.Context) *Session {
	session, err := GetSession(ctx)
	if err != nil {
		panic(err)
	}

	return session
}

// setSession sets the session in the context.
func setSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}
