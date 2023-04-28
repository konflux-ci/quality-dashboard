package login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var provider *oidc.Provider
var gh_config oauth2.Config
var verifier *oidc.IDTokenVerifier

const exampleAppState = "I wish to wash my irish wristwatch"

func init() {
	var err error
	provider, err = oidc.NewProvider(context.TODO(), "http://127.0.0.1:5556/dex")
	if err != nil {
		fmt.Println(err)
	}

	secret := os.Getenv("DEX_APP_SECRET")

	gh_config = oauth2.Config{
		ClientID:     "example-app",
		ClientSecret: secret,
		RedirectURL:  "http://localhost:9898/api/quality/login/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: "example-app"})
}

func (s *loginRouter) login(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

	// Configure an OpenID Connect aware OAuth2 client.

	authCodeURL := gh_config.AuthCodeURL(exampleAppState)
	fmt.Println(w.Header())

	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
	return nil
}

func (s *loginRouter) handleOAuth2Callback(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error { // Verify state and errors.

	var (
		err   error
		token *oauth2.Token
	)

	switch r.Method {
	case http.MethodGet:
		// Authorization redirect callback from OAuth2 auth flow.
		if errMsg := r.FormValue("error"); errMsg != "" {
			http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
			return err
		}
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, fmt.Sprintf("no code in request: %q", r.Form), http.StatusBadRequest)
			return err
		}
		if state := r.FormValue("state"); state != exampleAppState {
			http.Error(w, fmt.Sprintf("expected state %q got %q", exampleAppState, state), http.StatusBadRequest)
			return err
		}
		token, err = gh_config.Exchange(ctx, code)
	case http.MethodPost:
		// Form request from frontend to refresh a token.
		refresh := r.FormValue("refresh_token")
		if refresh == "" {
			http.Error(w, fmt.Sprintf("no refresh_token in request: %q", r.Form), http.StatusBadRequest)
			return err
		}
		t := &oauth2.Token{
			RefreshToken: refresh,
			Expiry:       time.Now().Add(-time.Hour),
		}
		token, err = gh_config.TokenSource(ctx, t).Token()
	default:
		http.Error(w, fmt.Sprintf("method not implemented: %s", r.Method), http.StatusBadRequest)
		return err
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get token: %v", err), http.StatusInternalServerError)
		return err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "no id_token in token response", http.StatusInternalServerError)
		return err
	}

	idToken, err := verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify ID token: %v", err), http.StatusInternalServerError)
		return err
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		http.Error(w, "no access_token in token response", http.StatusInternalServerError)
		return err
	}

	var claims json.RawMessage
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, fmt.Sprintf("error decoding ID token claims: %v", err), http.StatusInternalServerError)
		return err
	}

	buff := new(bytes.Buffer)
	if err := json.Indent(buff, []byte(claims), "", "  "); err != nil {
		http.Error(w, fmt.Sprintf("error indenting ID token claims: %v", err), http.StatusInternalServerError)
		return err
	}

	renderToken(w, gh_config.RedirectURL, rawIDToken, accessToken, token.RefreshToken, buff.String())
	return nil
}

var indexTmpl = template.Must(template.New("index.html").Parse(`<html>
  <head>
    <style>
form  { display: table;      }
p     { display: table-row;  }
label { display: table-cell; }
input { display: table-cell; }
    </style>
  </head>
  <body>
    <form action="/login" method="post">
      <p>
        <label> Authenticate for: </label>
        <input type="text" name="cross_client" placeholder="list of client-ids">
      </p>
      <p>
        <label>Extra scopes: </label>
        <input type="text" name="extra_scopes" placeholder="list of scopes">
      </p>
      <p>
        <label>Connector ID: </label>
        <input type="text" name="connector_id" placeholder="connector id">
      </p>
      <p>
        <label>Request offline access: </label>
        <input type="checkbox" name="offline_access" value="yes" checked>
      </p>
      <p>
	    <input type="submit" value="Login">
      </p>
    </form>
  </body>
</html>`))

func renderIndex(w http.ResponseWriter) {
	renderTemplate(w, indexTmpl, nil)
}

type tokenTmplData struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	RedirectURL  string
	Claims       string
}

var tokenTmpl = template.Must(template.New("token.html").Parse(`<html>
  <head>
    <style>
/* make pre wrap */
pre {
 white-space: pre-wrap;       /* css-3 */
 white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
 white-space: -pre-wrap;      /* Opera 4-6 */
 white-space: -o-pre-wrap;    /* Opera 7 */
 word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
    </style>
  </head>
  <body>
    <p> ID Token: <pre><code>{{ .IDToken }}</code></pre></p>
    <p> Access Token: <pre><code>{{ .AccessToken }}</code></pre></p>
    <p> Claims: <pre><code>{{ .Claims }}</code></pre></p>
	{{ if .RefreshToken }}
    <p> Refresh Token: <pre><code>{{ .RefreshToken }}</code></pre></p>
	<form action="{{ .RedirectURL }}" method="post">
	  <input type="hidden" name="refresh_token" value="{{ .RefreshToken }}">
	  <input type="submit" value="Redeem refresh token">
    </form>
	{{ end }}
  </body>
</html>
`))

func renderToken(w http.ResponseWriter, redirectURL, idToken, accessToken, refreshToken, claims string) {
	renderTemplate(w, tokenTmpl, tokenTmplData{
		IDToken:      idToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RedirectURL:  redirectURL,
		Claims:       claims,
	})
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *template.Error:
		// An ExecError guarantees that Execute has not written to the underlying reader.
		log.Printf("Error rendering template %s: %s", tmpl.Name(), err)

		// TODO(ericchiang): replace with better internal server error.
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		// An error with the underlying write, such as the connection being
		// dropped. Ignore for now.
	}
}
