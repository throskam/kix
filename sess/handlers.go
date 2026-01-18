package sess

import "net/http"

// Logout is a handler for logging out.
func Logout(
	store SessionStore,
	redirectURL string,
	handleError func(http.ResponseWriter, *http.Request, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := store.Erase(r, w)
		if err != nil {
			handleError(w, r, err)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
