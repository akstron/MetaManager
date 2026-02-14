package cmd

import (
	"context"
	"fmt"

	"github.com/heroku/self/MetaManager/internal/googleauth"
	"github.com/heroku/self/MetaManager/internal/services"
	"github.com/spf13/cobra"
)

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
	creds := services.EmbeddedCredentials()
	if len(creds) == 0 {
		return fmt.Errorf("no embedded credentials; rebuild the binary with credentials.json")
	}
	config, err := googleauth.LoadConfigFromBytes(creds)
	if err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}

	tok, err := googleauth.RunLoginFlow(config)
	if err != nil {
		return fmt.Errorf("login flow: %w", err)
	}

	tokenPath, err := services.TokenPath()
	if err != nil {
		return err
	}
	if err := googleauth.SaveToken(tokenPath, tok); err != nil {
		return fmt.Errorf("save token to %q: %w", tokenPath, err)
	}

	fmt.Printf("Token saved to %s\n", tokenPath)
	return nil
}

// GetGDriveService returns a GDrive service (uses credentials set in services via SetEmbeddedCredentials from main).
func GetGDriveService(ctx context.Context) (*services.GDriveService, error) {
	return services.GetGDriveService(ctx)
}
