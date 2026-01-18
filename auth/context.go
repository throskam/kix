package auth

import (
	"context"
	"errors"
)

// contextKey is a type for context keys.
type contextKey string

// identityContextKey is the context key for the identity.
const identityContextKey contextKey = "identity"

var (
	ErrCtxIdentityMissing = errors.New("missing identity context")
	ErrCtxSessionMissing  = errors.New("missing session context")
)

// GetIdentity returns the identity from the context.
func GetIdentity[T any](ctx context.Context) (T, error) {
	identity, ok := ctx.Value(identityContextKey).(T)
	if !ok {
		var empty T

		return empty, ErrCtxIdentityMissing
	}

	return identity, nil
}

// MustGetIdentity returns the identity from the context or panics if it is missing.
func MustGetIdentity[T any](ctx context.Context) T {
	identity, err := GetIdentity[T](ctx)
	if err != nil {
		panic(err)
	}

	return identity
}

// setIdentity sets the identity in the context.
func setIdentity(ctx context.Context, user any) context.Context {
	return context.WithValue(ctx, identityContextKey, user)
}
