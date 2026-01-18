package i18n

import (
	"context"

	"golang.org/x/text/message"
)

// contextKey is a type for context keys.
type contextKey string

// translatorContextKey is the context key for the translator.
const translatorContextKey contextKey = "translator"

// MustGetTranslator returns the translator from the context or panics if it is missing.
func MustGetTranslator(ctx context.Context) *message.Printer {
	translator, ok := ctx.Value(translatorContextKey).(*message.Printer)
	if !ok {
		panic("no translator")
	}

	return translator
}

// setTranslator sets the translator in the context.
func setTranslator(ctx context.Context, translator *message.Printer) context.Context {
	return context.WithValue(ctx, translatorContextKey, translator)
}
