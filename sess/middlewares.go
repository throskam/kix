package sess

import (
	"fmt"
	"net/http"

	"github.com/throskam/ki"
)

// Sessionizer is a middleware that adds a session to the request context.
func Sessionizer(
	store SessionStore,
	handleError func(http.ResponseWriter, *http.Request, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Read(r)
			if err != nil {
				err2 := store.Erase(r, w)
				if err2 != nil {
					handleError(w, r, fmt.Errorf("%w: %w", ErrSessionEraseFailure, err2))
					return
				}

				handleError(w, r, fmt.Errorf("%w: %w", ErrSessionReadFailure, err))
				return
			}

			ctx := setSession(r.Context(), session)

			brw := ki.NewBufferedResponseWriter(w)

			next.ServeHTTP(brw, r.WithContext(ctx))

			err = store.Write(r, brw, session)
			if err != nil {
				handleError(w, r, fmt.Errorf("%w: %w", ErrSessionWriteFailure, err))
				return
			}

			_, err = brw.Flush()
			if err != nil {
				handleError(w, r, fmt.Errorf("%w: %w", ErrSessionWriteFailure, err))
				return
			}
		})
	}
}
