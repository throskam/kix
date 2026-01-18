package auth

import (
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrJWTActiveKeyMissing   = errors.New("missing active key")
	ErrJWTClaimsParseFailure = errors.New("failed to parse claims")
	ErrJWTInvalid            = errors.New("invalid")
	ErrJWTKidClaimMissing    = errors.New("missing kid claim")
	ErrJWTKidClaimUnknown    = errors.New("unknown kid claim")
	ErrJWTSignFailure        = errors.New("failed to sign")
)

// JWK represents a JSON Web Key.
type JWK struct {
	Kid    string
	Value  []byte
	Active bool
}

// JWKS represents a map of JSON Web Keys.
type JWKS map[string]JWK

// GetByKid returns the JWK with the given key ID.
func (jwks JWKS) GetByKid(kid string) (JWK, bool) {
	jwk, ok := jwks[kid]

	return jwk, ok
}

// Add adds a JWK to the JWKS.
func (jwks JWKS) Add(jwk JWK) {
	jwks[jwk.Kid] = jwk
}

// Remove removes a JWK from the JWKS.
func (jwks JWKS) Remove(jwk JWK) {
	delete(jwks, jwk.Kid)
}

// GetActive returns the first active JWK or the first JWK if none are active.
func (jwks JWKS) GetActive() (JWK, bool) {
	for _, jwk := range jwks {
		if jwk.Active {
			return jwk, true
		}
	}

	for _, jwk := range jwks {
		return jwk, true
	}

	return JWK{}, false
}

// Claims represents the claims of a JWT.
type Claims map[string]any

// GenerateJWT generates a JWT.
func GenerateJWT(
	jwks JWKS,
	customClaims Claims,
	ttl time.Duration,
) (string, error) {
	jwk, ok := jwks.GetActive()
	if !ok {
		return "", ErrJWTActiveKeyMissing
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(ttl)),
		"nbf": jwt.NewNumericDate(time.Now()),
		"iat": jwt.NewNumericDate(time.Now()),
		"jti": uuid.NewString(),
	}

	maps.Copy(claims, customClaims)

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	token.Header["kid"] = jwk.Kid

	tokenString, err := token.SignedString(jwk.Value)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrJWTSignFailure, err)
	}

	return tokenString, nil
}

// ParseJWT parses a JWT and returns its claims.
func ParseJWT(
	jwks JWKS,
	bearer string,
	leeway time.Duration,
) (Claims, error) {
	mapClaims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(bearer, mapClaims, func(token *jwt.Token) (any, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, ErrJWTKidClaimMissing
		}

		jwk, ok := jwks.GetByKid(kid)
		if !ok {
			return nil, ErrJWTKidClaimUnknown
		}

		return jwk.Value, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithLeeway(leeway),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("%w: %w", ErrJWTClaimsParseFailure, err)
	}

	if !token.Valid {
		return Claims{}, ErrJWTInvalid
	}

	claims := Claims{}

	maps.Copy(claims, *mapClaims)

	return claims, nil
}
