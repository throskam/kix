package i18n

import (
	"context"

	"golang.org/x/text/message"
)

// T returns the translated string for the given key and arguments.
func T(ctx context.Context, key message.Reference, args ...any) string {
	return MustGetTranslator(ctx).Sprintf(key, args...)
}
