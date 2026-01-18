package i18n

import (
	"net/http"
	"sync"

	"github.com/throskam/ki"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Translator returns a middleware that loads the translator for the request language.
func Translator() func(http.Handler) http.Handler {
	var translators sync.Map

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := ki.MustGetLanguage(r.Context())

			translator := loadOrStoreTranslator(
				&translators,
				lang,
			)

			ctx := setTranslator(r.Context(), translator)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func loadOrStoreTranslator(
	translators *sync.Map,
	lang language.Tag,
) *message.Printer {
	key := lang.String()

	if translator, ok := translators.Load(key); ok {
		return translator.(*message.Printer)
	}

	printer := message.NewPrinter(lang)

	translators.Store(key, printer)

	return printer
}
