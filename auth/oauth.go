package auth

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/throskam/kix/sess"
	"golang.org/x/oauth2"
)

var (
	ErrOAuthAuthenticateFailure = errors.New("failed to authenticate")
	ErrOAuthCSRFTokenMismatch   = errors.New("CSRF mismatch")
	ErrOAuthExchangeFailure     = errors.New("failed to exchange authorization code")
	ErrOAuthInitiateFailure     = errors.New("failed to initiate")
)

// OAuthCSRFTokenKey is the session key for the OAuth CSRF token.
var OAuthCSRFTokenKey = "oauth-csrf-token"

// OAuth2AuthStrategy is the interface for OAuth2 authentication strategies.
type OAuth2AuthStrategy interface {
	Initiate(*http.Request, *sess.Session) error
	Authenticate(context.Context, *oauth2.Token, *sess.Session) (string, error)
	HandleError(http.ResponseWriter, *http.Request, error)
}

// OAuth2Controller is a controller for OAuth2 authentication.
type OAuth2Controller struct {
	config   *oauth2.Config
	strategy OAuth2AuthStrategy
}

// NewOAuth2Controller creates a new OAuth2Controller.
func NewOAuth2Controller(
	config *oauth2.Config,
	strategy OAuth2AuthStrategy,
) *OAuth2Controller {
	return &OAuth2Controller{
		config:   config,
		strategy: strategy,
	}
}

// Login initiates the OAuth2 login flow.
func (p *OAuth2Controller) Login(w http.ResponseWriter, r *http.Request) {
	session := sess.MustGetSession(r.Context())

	err := p.strategy.Initiate(r, session)
	if err != nil {
		p.strategy.HandleError(w, r, fmt.Errorf("%w: %w", ErrOAuthInitiateFailure, err))
		return
	}

	// CSRF token is used here to protect against cross-site request forgery attacks
	// during the OAuth2 login flow. It ensures that the OAuth2 callback corresponds
	// to the login request initiated by the same user. By storing a unique token in
	// the user's session and verifying it upon callback, we can confirm that the
	// response from the OAuth provider was not forged or initiated by a third party.
	csrfToken := GenerateCSRFToken()
	providerURL := p.config.AuthCodeURL(csrfToken)
	session.Add(OAuthCSRFTokenKey, csrfToken)

	http.Redirect(w, r, providerURL, http.StatusSeeOther)
}

// Callback handles the OAuth2 callback.
func (p *OAuth2Controller) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	// Exchange code for an access token.
	token, err := p.config.Exchange(r.Context(), code)
	if err != nil {
		p.strategy.HandleError(w, r, fmt.Errorf("%w: %w", ErrOAuthExchangeFailure, err))
		return
	}

	session := sess.MustGetSession(r.Context())
	csrfToken := session.GetFirst(OAuthCSRFTokenKey)
	session.Del(OAuthCSRFTokenKey)

	// Check that the CSRF token created during the login matches the one we receive from the callback.
	if csrfToken != r.FormValue("state") {
		p.strategy.HandleError(w, r, ErrOAuthCSRFTokenMismatch)
		return
	}

	redirectURL, err := p.strategy.Authenticate(r.Context(), token, session)
	if err != nil {
		p.strategy.HandleError(w, r, fmt.Errorf("%w: %w", ErrOAuthAuthenticateFailure, err))
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func init() {
	gob.Register(url.Values{})
}
