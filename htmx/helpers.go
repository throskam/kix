package htmx

import "net/http"

// Redirect redirects the user to the given URL using HTMX headers.
func Redirect(w http.ResponseWriter, r *http.Request, redirectURL string) {
	if r.Header.Get("HX-Request") == "true" {
		if r.Header.Get("HX-Boosted") == "true" {
			w.Header().Add("HX-Location", redirectURL)
			w.WriteHeader(200)
			return
		}

		// Scenario:
		// we are in the middle of fetching a section of the page and we want to
		// trigger a redirection. This should never happen and is probably due
		// to permissions issues so we force a page refresh and hope for the best.
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(200)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
