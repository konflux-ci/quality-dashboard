package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
)

// Authentication verify the athenticity of the tokens spawned by dex.
type Authentication struct {
	IssuerUrl           string
	ApplicationClientId string
}

type Claims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"email_verified"`
	Groups   []string `json:"groups"`
}

// NewAuthenticationMiddleware creates a new Authentication protect backend routes.
func NewAuthenticationMiddleware(url string, clientId string) Authentication {
	return Authentication{IssuerUrl: url,
		ApplicationClientId: clientId,
	}
}

// WrapHandler returns a new handler function wrapping the previous one in the request chain.
func (c Authentication) WrapHandler(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error) func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		token := ExtractBearerToken(r.Header.Get("Authorization"))

		if token == "" {
			return fmt.Errorf("token not present in headers")
		}

		if err := c.Authorize(context.TODO(), token); err != nil {
			return fmt.Errorf("failed to verify token, %s", err)
		}

		return handler(ctx, w, r, vars)
	}
}

// extractBearerToken extracts the Bearer token from the Authorization header
func ExtractBearerToken(token string) string {
	return strings.TrimSpace(strings.Replace(token, "Bearer ", "", 1))
}

// Authorize verifies a bearer token and pulls user information form the claims.
func (c Authentication) Authorize(ctx context.Context, bearerToken string) error {
	var claims Claims
	// Initialize a provider by specifying dex's issuer URL.
	provider, err := oidc.NewProvider(ctx, c.IssuerUrl)
	if err != nil {
		return fmt.Errorf("error fetching provider: %s", err)
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: c.ApplicationClientId})

	idToken, err := idTokenVerifier.Verify(ctx, bearerToken)
	if err != nil {
		return fmt.Errorf("could not verify bearer token: %v", err)
	}

	if err := idToken.Claims(&claims); err != nil {
		return fmt.Errorf("failed to parse claims: %v", err)
	}

	if !claims.Verified {
		return fmt.Errorf("email (%q) in returned claims was not verified", claims.Email)
	}

	return nil
}
