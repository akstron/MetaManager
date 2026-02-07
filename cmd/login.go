package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/googleauth"
	"github.com/spf13/cobra"
)

const (
	googleTokenFileName = "google_token.json"
)

// EmbeddedCredentials is set by main when credentials.json is embedded in the binary.
var embeddedCredentials []byte

// SetEmbeddedCredentials sets the credentials JSON embedded in the binary (called from main).
func SetEmbeddedCredentials(b []byte) {
	embeddedCredentials = b
}

// loginCmd runs Google OAuth (desktop flow) and stores the token next to the binary.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in with Google and store the token for Drive access",
	Long:  `Uses embedded credentials to run the OAuth flow, then saves the token to google_token.json in the same directory as the installed binary.`,
	RunE:  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	if len(embeddedCredentials) == 0 {
		return fmt.Errorf("no embedded credentials; rebuild the binary with credentials.json")
	}
	config, err := googleauth.LoadConfigFromBytes(embeddedCredentials)
	if err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}

	tok, err := googleauth.RunLoginFlow(config)
	if err != nil {
		return fmt.Errorf("login flow: %w", err)
	}

	tokenPath, err := resolveTokenPath()
	if err != nil {
		return err
	}
	if err := googleauth.SaveToken(tokenPath, tok); err != nil {
		return fmt.Errorf("save token to %q: %w", tokenPath, err)
	}

	fmt.Printf("Token saved to %s\n", tokenPath)
	return nil
}

// resolveTokenPath returns the path to google_token.json in the same directory as the executable.
func resolveTokenPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, googleTokenFileName), nil
}
