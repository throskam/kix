package htmx

import (
	"net/http"
	"net/url"

	"github.com/throskam/ki"
	"golang.org/x/text/message"
)

// FormData represents the payload of a request.
type FormData interface {
	Parse(*http.Request)
	Validate(*message.Printer) url.Values
}

// Form represents a form.
type Form[T FormData] struct {
	Data       T
	Validation url.Values
}

// NewFormFromRequest creates a new form from the given request.
func NewFormFromRequest[T FormData](r *http.Request, data T) *Form[T] {
	data.Parse(r)

	p := message.NewPrinter(ki.MustGetLanguage(r.Context()))

	return &Form[T]{
		Data:       data,
		Validation: data.Validate(p),
	}
}

// NewForm creates a new form.
func NewForm[T FormData](data T) *Form[T] {
	return &Form[T]{
		Data:       data,
		Validation: url.Values{},
	}
}

// OK returns true if the form is valid.
func (f Form[T]) OK() bool {
	for _, messages := range f.Validation {
		if len(messages) > 0 {
			return false
		}
	}

	return true
}
