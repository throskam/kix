package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/throskam/kix/sess"
)

var ErrAuthIdentifyFailure = errors.New("failed to identify")

// Authenticate authenticates the user and sets the identity in the context.
func Authenticate(
	identify func(context.Context, *sess.Session) (any, error),
	handleError func(http.ResponseWriter, *http.Request, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := sess.MustGetSession(r.Context())

			identity, err := identify(r.Context(), session)
			if err != nil {
				handleError(w, r, fmt.Errorf("%w: %w", ErrAuthIdentifyFailure, err))
				return
			}

			ctx := setIdentity(r.Context(), identity)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Authenticated ensures that the user is authenticated.
func Authenticated[T any](
	handleError func(http.ResponseWriter, *http.Request, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetIdentity[T](r.Context())
			if err != nil {
				handleError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Anonymous[T any](
	handleError func(http.ResponseWriter, *http.Request, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetIdentity[T](r.Context())
			if err == nil {
				handleError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRF protects against cross-site request forgery attacks.
func CSRF(
	handleError func(http.ResponseWriter, *http.Request, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// GET requests are exempt from CSRF protection.
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			token := GetCSRFToken(r.Context())

			if token == "" || token != r.Header.Get("X-CSRF-TOKEN") {
				handleError(w, r, ErrCSRFTokenMismatch)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
