// Package googleauth provides desktop OAuth flow for Google (e.g. Drive) using
// credentials.json and stores the token for reuse.
package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	driveScope   = "https://www.googleapis.com/auth/drive.metadata.readonly"
	callbackPath = "/callback"
	callbackPort = "8080"
)

// LoadConfigFromFile reads desktop or web client credentials from a local JSON file
// and returns an oauth2.Config with redirect URL set to localhost for the login flow.
func LoadConfigFromFile(path string) (*oauth2.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read credentials file: %w", err)
	}
	return LoadConfigFromBytes(b)
}

// LoadConfigFromBytes builds an oauth2.Config from desktop/web client credentials JSON
// (e.g. embedded in the binary). Redirect URL is set to localhost for the login flow.
func LoadConfigFromBytes(b []byte) (*oauth2.Config, error) {
	config, err := google.ConfigFromJSON(b, driveScope)
	if err != nil {
		return nil, fmt.Errorf("parse credentials JSON: %w", err)
	}
	config.RedirectURL = fmt.Sprintf("http://localhost:%s%s", callbackPort, callbackPath)
	return config, nil
}

var (
	stateStore   = make(map[string]time.Time)
	stateStoreMu sync.RWMutex
)

func addState() string {
	s := fmt.Sprintf("%d-%s", time.Now().UnixNano(), mustRandomHex(16))
	stateStoreMu.Lock()
	defer stateStoreMu.Unlock()
	stateStore[s] = time.Now().Add(10 * time.Minute)
	return s
}

func validateState(s string) bool {
	stateStoreMu.Lock()
	defer stateStoreMu.Unlock()
	exp, ok := stateStore[s]
	if !ok || time.Now().After(exp) {
		return false
	}
	delete(stateStore, s)
	return true
}

func mustRandomHex(n int) string {
	b := make([]byte, n)
	f, err := os.Open("/dev/urandom")
	if err != nil {
		for i := range b {
			b[i] = byte(time.Now().UnixNano() % 256)
		}
		return fmt.Sprintf("%x", b)
	}
	defer f.Close()
	_, _ = f.Read(b)
	return fmt.Sprintf("%x", b)
}

// RunLoginFlow runs the OAuth flow: starts a local server, opens the browser for
// user sign-in, and returns the token after the callback.
func RunLoginFlow(config *oauth2.Config) (*oauth2.Token, error) {
	state := addState()
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("invalid state parameter")
			http.Error(w, "Authentication failed: invalid state.", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("missing code in callback")
			http.Error(w, "Authentication failed: no code received.", http.StatusBadRequest)
			return
		}
		codeCh <- code
		_, _ = w.Write([]byte(`<html><body><h2>Success!</h2><p>You can close this tab and return to the terminal.</p></body></html>`))
	})

	srv := &http.Server{Addr: ":" + callbackPort, Handler: mux}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Opening browser for Google sign-in. If it doesn't open, visit:\n  %s\n", authURL)
	_ = browser.OpenURL(authURL)

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		_ = srv.Shutdown(context.Background())
		wg.Wait()
		return nil, err
	}

	_ = srv.Shutdown(context.Background())
	wg.Wait()

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("exchange code for token: %w", err)
	}
	return tok, nil
}

// SaveToken writes the token as JSON to path for future use.
func SaveToken(path string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadToken reads a previously saved token from path.
func LoadToken(path string) (*oauth2.Token, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tok oauth2.Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}
